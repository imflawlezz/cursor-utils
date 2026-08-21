package github

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
)

func TestSafeTag(t *testing.T) {
	t.Parallel()
	if !SafeTag("v0.1.0") || !SafeTag("0.1.0-rc.1") {
		t.Fatal("valid tags rejected")
	}
	for _, tag := range []string{"", "../v1", "v1/foo", "foo bar", "v1?ref=x", "v1\\x", string([]byte{'v', 0, '1'})} {
		if SafeTag(tag) {
			t.Errorf("unsafe tag accepted: %q", tag)
		}
	}
}

func TestExtractTarGzSkipsNonComponentsAndRejectsTraversal(t *testing.T) {
	t.Parallel()
	ok := pack(t, map[string]string{
		"cursor-utils-v0.1.0/README.md":          "# hi",
		"cursor-utils-v0.1.0/commands/commit.md": "commit",
		"cursor-utils-v0.1.0/LICENSE":            "mit",
	})
	bundle, err := ExtractTarGz(bytes.NewReader(ok), "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if !bundle.HasComponent("commands") || len(bundle.AllFiles()) != 1 {
		t.Fatalf("files = %+v", bundle.AllFiles())
	}

	evil := pack(t, map[string]string{
		"cursor-utils-v0.1.0/commands/../../etc/passwd": "nope",
	})
	if _, err := ExtractTarGz(bytes.NewReader(evil), "v0.1.0"); err == nil {
		t.Fatal("expected traversal to be rejected")
	}
}

func TestExtractTarGzSkipsSymlinks(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	body := []byte("ok")
	hdr := &tar.Header{Name: "repo/commands/commit.md", Mode: 0644, Size: int64(len(body)), Typeflag: tar.TypeReg}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatal(err)
	}
	link := &tar.Header{Name: "repo/commands/link.md", Mode: 0644, Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd"}
	if err := tw.WriteHeader(link); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	bundle, err := ExtractTarGz(bytes.NewReader(buf.Bytes()), "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.AllFiles()) != 1 {
		t.Fatalf("got %d files", len(bundle.AllFiles()))
	}
}

func TestClientListTagsAndFetch(t *testing.T) {
	t.Parallel()
	archive := pack(t, map[string]string{
		"cursor-utils-v0.1.0/commands/docs.md": "docs",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/imflawlezz/cursor-utils/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `[{"name":"v0.1.0"},{"name":"../evil"},{"name":"v0.2.0"}]`)
	})
	mux.HandleFunc("/imflawlezz/cursor-utils/archive/refs/tags/v0.1.0.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := New(config.Default())
	c.HTTP = srv.Client()
	c.APIURL = srv.URL
	c.WebURL = srv.URL

	tags, err := c.ListTags(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 || tags[0] != "v0.1.0" || tags[1] != "v0.2.0" {
		t.Fatalf("tags = %#v", tags)
	}

	bundle, err := c.Fetch(context.Background(), "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := bundle.File("commands/docs.md"); !ok {
		t.Fatal("missing file")
	}
	if _, err := c.Fetch(context.Background(), "../evil"); err == nil {
		t.Fatal("expected unsafe tag to fail")
	}
}

func TestClientHTTPErrors(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(srv.Close)
	c := New(config.Default())
	c.HTTP = srv.Client()
	c.APIURL = srv.URL
	_, err := c.ListTags(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*HTTPError); !ok {
		t.Fatalf("got %T %v", err, err)
	}
}

func pack(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, body := range files {
		b := []byte(body)
		hdr := &tar.Header{Name: name, Mode: 0644, Size: int64(len(b)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(b); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
