//go:build !windows

package store

import "os"

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// replace puts tmp in path's place.
func replace(tmp, path string) error {
	return os.Rename(tmp, path)
}
