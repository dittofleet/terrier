package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/uninstall"
	"github.com/dittofleet/go-cli-kit/xdg"
	"github.com/dittofleet/terrier/internal/app"
	"github.com/dittofleet/terrier/internal/store"
)

const uninstallUsage = "usage: terrier uninstall [--yes]"

// Uninstall removes the terrier binary, its `ter` alias, and the config
// directory.
func Uninstall(args []string, a clikit.App) error {
	var yes bool
	rest, err := yesFlag(&yes).parse(args, uninstallUsage)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("unexpected argument: %s\n%s", rest[0], uninstallUsage)
	}

	items := []uninstall.Item{
		{Label: "Config", Path: xdg.ConfigDir(app.Name), Note: describeRegistry(), Remove: os.RemoveAll},
	}
	// If the binary cannot be found, the kit says so before removing anything.
	if binary, err := clikit.Executable(); err == nil {
		if alias := aliasPath(binary); alias != "" {
			items = append(items, uninstall.Item{Label: "Alias", Path: alias, Remove: os.Remove})
		}
	}
	return uninstall.Run(a, yes, uninstall.Plan{
		Items: items,
		Notice: "Only the registry is deleted. No repository is touched, and tools that read\n" +
			"terrier keep whatever they stored themselves.",
	})
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

// describeRegistry summarizes what is about to be deleted. An unreadable
// registry is not worth failing the uninstall over.
func describeRegistry() string {
	s, err := store.Load()
	if err != nil {
		return "your registered projects"
	}
	return plural(len(s.Projects), "project")
}
