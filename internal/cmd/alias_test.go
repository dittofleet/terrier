package cmd

import (
	"os"
	"path/filepath"
	"testing"

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

	alias := filepath.Join(dir, aliasFile)
	if err := createAlias(binary, alias); err != nil {
		t.Fatal(err)
	}
	if got := aliasPath(binary); got != alias {
		t.Errorf("got %q, want %q", got, alias)
	}
}

// Anything else at that name belongs to something else, and uninstall must
// not take it for terrier's.
func TestAliasIgnoresAForeignFile(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "terrier")
	if err := os.WriteFile(filepath.Join(dir, aliasFile), []byte("someone else's"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := aliasPath(binary); got != "" {
		t.Errorf("a foreign file was taken for the alias: %q", got)
	}
}
