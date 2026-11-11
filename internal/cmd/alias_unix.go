//go:build !windows

package cmd

import (
	"os"
	"path/filepath"

	"github.com/dittofleet/terrier/internal/app"
)

// aliasFile is the alias's name on disk.
const aliasFile = app.Alias

// createAlias symlinks alias to the binary beside it. A symlink rather
// than a copy, so `terrier update` moves both at once.
func createAlias(binary, alias string) error {
	return os.Symlink(filepath.Base(binary), alias)
}

// isAlias reports whether alias is terrier's own link to binary.
func isAlias(alias, binary string) bool {
	resolved, err := filepath.EvalSymlinks(alias)
	return err == nil && resolved == binary
}
