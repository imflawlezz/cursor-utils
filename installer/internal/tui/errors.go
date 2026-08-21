package tui

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/imflawlezz/cursor-utils/installer/internal/github"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
)

func friendly(err error) string {
	if err == nil {
		return ""
	}
	var fe *manifest.FormatError
	if errors.As(err, &fe) {
		return fe.Error() + "\n\nThe installer will not change Cursor files until the manifest can be read safely."
	}
	var he *github.HTTPError
	if errors.As(err, &he) {
		switch he.Status {
		case 403, 429:
			return "GitHub rate limit reached.\n\nWait a few minutes and try again."
		case 404:
			return "The requested version was not found on GitHub.\n\nChoose another version."
		default:
			if he.Status >= 500 {
				return "GitHub is temporarily unavailable.\n\nTry again later."
			}
		}
	}
	if errors.Is(err, context.Canceled) {
		return "The operation was cancelled."
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "The request timed out.\n\nCheck your network connection and try again."
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return "The request to GitHub timed out.\n\nCheck your network connection and try again."
	}
	msg := err.Error()
	if strings.Contains(strings.ToLower(msg), "unable to reach github") ||
		strings.Contains(strings.ToLower(msg), "no such host") ||
		strings.Contains(strings.ToLower(msg), "connection refused") {
		return "Unable to reach GitHub.\n\nCheck your network connection and try again."
	}
	return fmt.Sprintf("%s", msg)
}
