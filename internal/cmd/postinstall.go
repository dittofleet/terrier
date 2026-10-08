package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/postinstall"
	"github.com/dittofleet/terrier/internal/app"
)

// Postinstall is the install script's first-time setup: the `ter` alias
// beside the binary. A symlink rather than a copy, so `terrier update`
// moves both at once. Anything already at that name belongs to something
// else and is left alone.
func Postinstall(a clikit.App) error {
	return postinstall.Run(a, func() error {
		binary, err := clikit.Executable()
		if err != nil {
			return err
		}
		if aliasPath(binary) != "" {
			return nil
		}
		alias := filepath.Join(filepath.Dir(binary), app.Alias)
		if _, err := os.Lstat(alias); err == nil {
			fmt.Printf("Note: %s already exists and is not terrier's alias. Leaving it alone.\n", alias)
			return nil
		}
		return os.Symlink(filepath.Base(binary), alias)
	})
}
