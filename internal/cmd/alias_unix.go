//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dittofleet/terrier/internal/app"
)

// setupAlias symlinks `ter` to the binary beside it. A symlink rather than
// a copy, so `terrier update` moves both at once. Anything already at that
// name belongs to something else and is left alone.
func setupAlias(binary string) error {
	if aliasPath(binary) != "" {
		return nil
	}
	alias := filepath.Join(filepath.Dir(binary), app.Alias)
	if _, err := os.Lstat(alias); err == nil {
		fmt.Printf("Note: %s already exists and is not terrier's alias. Leaving it alone.\n", alias)
		return nil
	}
	return os.Symlink(filepath.Base(binary), alias)
}

// aliasPath returns the `ter` symlink installed beside the binary, or ""
// when there is none there to remove. Anything at that name which is not
// a symlink to this binary belongs to something else and is left alone.
func aliasPath(binaryPath string) string {
	candidate := filepath.Join(filepath.Dir(binaryPath), app.Alias)
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil || resolved != binaryPath {
		return ""
	}
	return candidate
}
