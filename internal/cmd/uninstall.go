package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/xdg"
	"github.com/dittofleet/terrier/internal/app"
	"github.com/dittofleet/terrier/internal/store"
)

const uninstallUsage = "usage: terrier uninstall [--yes]"

// Uninstall removes the terrier binary, its `ter` alias, and the config
// directory. Order is config, then alias, then binary, so a failure leaves
// a tool to retry with.
func Uninstall(args []string, a clikit.App) error {
	var yes bool
	rest, err := yesFlag(&yes).parse(args, uninstallUsage)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("unexpected argument: %s\n%s", rest[0], uninstallUsage)
	}

	if a.IsDev() {
		return errors.New("cannot uninstall a dev build")
	}

	binaryPath, err := clikit.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine binary path: %w", err)
	}
	alias := aliasPath(binaryPath)

	configDir := xdg.ConfigDir(app.Name)

	fmt.Println("This will remove:")
	fmt.Printf("  - Binary:  %s\n", binaryPath)
	if alias != "" {
		fmt.Printf("  - Alias:   %s\n", alias)
	}
	fmt.Printf("  - Config:  %s  (%s)\n", configDir, describeRegistry())
	fmt.Println()
	fmt.Println("Only the registry is deleted. No repository is touched, and tools that read")
	fmt.Println("terrier keep whatever they stored themselves.")
	fmt.Println()

	if !yes {
		ok, err := confirm("Proceed?")
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Aborted.")
			return nil
		}
	}

	steps := []struct {
		label string
		path  string
		fn    func(string) error
	}{
		{"config directory", configDir, os.RemoveAll},
		{"alias", alias, os.Remove},
		{"binary", binaryPath, os.Remove},
	}
	var removed []string
	for _, s := range steps {
		if s.path == "" {
			continue
		}
		err := s.fn(s.path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			if len(removed) > 0 {
				fmt.Fprintf(os.Stderr, "Removed before failure: %s\n", strings.Join(removed, ", "))
			}
			return fmt.Errorf("failed to remove %s (%s): %w", s.label, s.path, err)
		}
		removed = append(removed, s.label)
	}

	fmt.Println("Uninstalled terrier.")
	return nil
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
