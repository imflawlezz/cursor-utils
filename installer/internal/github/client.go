package github

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
)

const (
	defaultAPI = "https://api.github.com"
	defaultWeb = "https://github.com"

	downloadTimeout = 60 * time.Second
	maxResponse     = 8 << 20
	maxArchive      = 20 << 20
)

type Client struct {
	HTTP   *http.Client
	APIURL string
	WebURL string
	Owner  string
	Repo   string
	UA     string
}

func New(cfg config.App) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: downloadTimeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 20 * time.Second,
				IdleConnTimeout:       30 * time.Second,
			},
		},
		APIURL: defaultAPI,
		WebURL: defaultWeb,
		Owner:  cfg.GitHubOwner,
		Repo:   cfg.GitHubRepo,
		UA:     "cursor-utils-installer/" + cfg.Version,
	}
}

type tagDTO struct {
	Name string `json:"name"`
}

func (c *Client) ListTags(ctx context.Context) ([]string, error) {
	var tags []string
	endpoint := strings.TrimRight(c.APIURL, "/") + "/repos/" + url.PathEscape(c.Owner) + "/" + url.PathEscape(c.Repo) + "/tags?per_page=100"
	for endpoint != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		c.headers(req)
		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		body, err := readLimited(resp.Body, maxResponse)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if err := checkStatus(resp, body); err != nil {
			return nil, err
		}
		var page []tagDTO
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("GitHub returned an unexpected tags response")
		}
		for _, t := range page {
			if SafeTag(t.Name) {
				tags = append(tags, t.Name)
			}
		}
		endpoint = nextLink(resp.Header.Get("Link"))
	}
	return tags, nil
}

func (c *Client) Fetch(ctx context.Context, tag string) (*content.Bundle, error) {
	if !SafeTag(tag) {
		return nil, fmt.Errorf("invalid version tag %q", tag)
	}
	archive := strings.TrimRight(c.WebURL, "/") + "/" + url.PathEscape(c.Owner) + "/" + url.PathEscape(c.Repo) + "/archive/refs/tags/" + url.PathEscape(tag) + ".tar.gz"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, archive, nil)
	if err != nil {
		return nil, err
	}
	c.headers(req)
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := readLimited(resp.Body, 4096)
		return nil, checkStatus(resp, body)
	}
	limited := io.LimitReader(resp.Body, maxArchive+1)
	bundle, err := ExtractTarGz(limited, tag)
	if err != nil {
		return nil, err
	}
	return bundle, nil
}

func (c *Client) headers(req *http.Request) {
	req.Header.Set("User-Agent", c.UA)
	req.Header.Set("Accept", "application/vnd.github+json, application/json, application/octet-stream")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, wrapNet(err)
	}
	return resp, nil
}

// SafeTag is true only for tags that can be interpolated into a GitHub URL path.
func SafeTag(tag string) bool {
	if tag == "" || len(tag) > 128 {
		return false
	}
	if strings.Contains(tag, "..") {
		return false
	}
	if strings.ContainsAny(tag, `/\?&#% `) {
		return false
	}
	for _, r := range tag {
		if unicode.IsControl(r) {
			return false
		}
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-' || r == '+') {
			return false
		}
	}
	return true
}

func ExtractTarGz(r io.Reader, tag string) (*content.Bundle, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("content archive is not valid gzip")
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	bundle := content.NewBundle(tag)
	var total int64
	nfiles := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("content archive is corrupt")
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		rel, ok := stripTop(hdr.Name)
		if !ok {
			continue
		}
		rel = strings.ReplaceAll(rel, `\`, "/")
		for _, part := range strings.Split(rel, "/") {
			if part == ".." || part == "." {
				return nil, fmt.Errorf("rejected unsafe path %q from remote content", rel)
			}
		}
		if strings.HasPrefix(rel, "/") {
			return nil, fmt.Errorf("rejected unsafe path %q from remote content", rel)
		}
		rel = path.Clean(rel)
		comp := content.ComponentOf(rel)
		if !content.Known(comp) {
			continue
		}
		if err := content.ValidateRelPath(rel); err != nil {
			return nil, fmt.Errorf("rejected unsafe path %q from remote content", rel)
		}
		if hdr.Size < 0 || hdr.Size > content.MaxFileSize {
			return nil, fmt.Errorf("%s is larger than the allowed file size", rel)
		}
		data, err := io.ReadAll(io.LimitReader(tr, hdr.Size+1))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s from the archive", rel)
		}
		if int64(len(data)) != hdr.Size {
			return nil, fmt.Errorf("archive entry size mismatch for %s", rel)
		}
		nfiles++
		if nfiles > content.MaxBundleFiles {
			return nil, fmt.Errorf("archive contains too many files")
		}
		total += int64(len(data))
		if total > content.MaxBundleBytes {
			return nil, fmt.Errorf("archive is larger than the allowed size")
		}
		if err := bundle.Add(rel, data); err != nil {
			return nil, err
		}
	}
	if bundle.Empty() {
		return nil, fmt.Errorf("version %s contains no installable Cursor content", tag)
	}
	return bundle, nil
}

func stripTop(name string) (string, bool) {
	name = strings.ReplaceAll(name, `\`, "/")
	name = strings.TrimPrefix(name, "./")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) < 2 || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func nextLink(header string) string {
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, `rel="next"`) {
			continue
		}
		start := strings.Index(part, "<")
		end := strings.Index(part, ">")
		if start >= 0 && end > start {
			return part[start+1 : end]
		}
	}
	return ""
}

func readLimited(r io.Reader, n int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, n+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > n {
		return nil, fmt.Errorf("GitHub response is too large")
	}
	return data, nil
}

type HTTPError struct {
	Status int
	Body   string
}

func (e *HTTPError) Error() string {
	switch e.Status {
	case http.StatusNotFound:
		return "the requested GitHub resource was not found"
	case http.StatusForbidden, http.StatusTooManyRequests:
		return "GitHub rate limit reached"
	default:
		if e.Status >= 500 {
			return "GitHub is temporarily unavailable"
		}
		return fmt.Sprintf("GitHub returned HTTP %d", e.Status)
	}
}

func checkStatus(resp *http.Response, body []byte) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return &HTTPError{Status: resp.StatusCode, Body: string(body)}
}

func wrapNet(err error) error {
	if err == nil {
		return nil
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return fmt.Errorf("the request to GitHub timed out: %w", err)
	}
	return fmt.Errorf("unable to reach GitHub: %w", err)
}
