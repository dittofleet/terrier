package cmd

import (
	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/postinstall"
)

// Postinstall is the install script's first-time setup: the `ter` alias
// beside the binary.
func Postinstall(a clikit.App) error {
	return postinstall.Run(a, func() error {
		binary, err := clikit.Executable()
		if err != nil {
			return err
		}
		return setupAlias(binary)
	})
}
