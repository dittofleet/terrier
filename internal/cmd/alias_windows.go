package cmd

import (
	"os"
	"path/filepath"

	"github.com/dittofleet/terrier/internal/app"
)

// aliasFile is the alias's name on disk. Windows finds a .cmd by its bare
// name, through PATHEXT.
const aliasFile = app.Alias + ".cmd"

// createAlias writes a batch file that runs the binary beside it. Creating
// a symlink on Windows takes Developer Mode or an administrator, and a
// hard link or a copy would keep running the old binary after an update.
func createAlias(binary, alias string) error {
	return os.WriteFile(alias, []byte(aliasScript(binary)), 0o755)
}

// isAlias reports whether alias is the batch file createAlias writes.
func isAlias(alias, binary string) bool {
	data, err := os.ReadFile(alias)
	return err == nil && string(data) == aliasScript(binary)
}

func aliasScript(binary string) string {
	return `@"%~dp0` + filepath.Base(binary) + `" %*` + "\r\n"
}
