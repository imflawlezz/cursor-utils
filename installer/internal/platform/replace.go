package platform

import (
	"os"
	"runtime"
)

// ReplaceFile moves src onto dest. Windows cannot os.Rename over an existing dest.
func ReplaceFile(src, dest string) error {
	if runtime.GOOS == "windows" {
		if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return os.Rename(src, dest)
}
