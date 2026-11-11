//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dittofleet/terrier/internal/app"
	"github.com/dittofleet/terrier/internal/git"
)

func TestAliasIsFoundOnceCreated(t *testing.T) {
	dir := git.Canonical(t.TempDir())
	binary := filepath.Join(dir, "terrier")
	if err := os.WriteFile(binary, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := aliasPath(binary); got != "" {
		t.Errorf("found an alias before creating one: %q", got)
	}

	if err := setupAlias(binary); err != nil {
		t.Fatal(err)
	}
	if got, want := aliasPath(binary), filepath.Join(dir, app.Alias); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Anything else at that name belongs to something else, and uninstall must
// not take it for terrier's.
func TestAliasLeavesAForeignFileAlone(t *testing.T) {
	dir := git.Canonical(t.TempDir())
	binary := filepath.Join(dir, "terrier")
	foreign := filepath.Join(dir, app.Alias)
	if err := os.WriteFile(foreign, []byte("someone else's"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := setupAlias(binary); err != nil {
		t.Fatal(err)
	}
	if got := aliasPath(binary); got != "" {
		t.Errorf("a foreign file was taken for the alias: %q", got)
	}
	if data, err := os.ReadFile(foreign); err != nil || string(data) != "someone else's" {
		t.Errorf("the foreign file changed: %q, %v", data, err)
	}
}
