package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/postinstall"
)

// Postinstall is the install script's first-time setup: the `ter` alias
// beside the binary. Anything already at that name belongs to something
// else and is left alone.
func Postinstall(a clikit.App) error {
	return postinstall.Run(a, func() error {
		binary, err := clikit.Executable()
		if err != nil {
			return err
		}
		alias := filepath.Join(filepath.Dir(binary), aliasFile)
		if isAlias(alias, binary) {
			return nil
		}
		if _, err := os.Lstat(alias); err == nil {
			fmt.Printf("Note: %s already exists and is not terrier's alias. Leaving it alone.\n", alias)
			return nil
		}
		return createAlias(binary, alias)
	})
}
