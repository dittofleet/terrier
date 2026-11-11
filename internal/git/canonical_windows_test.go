package git

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCanonicalResolvesAJunction(t *testing.T) {
	base := Canonical(t.TempDir())
	target := filepath.Join(base, "target")
	mkdirAll(t, target)
	junction := filepath.Join(base, "junction")
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", junction, target).CombinedOutput(); err != nil {
		t.Fatalf("mklink: %v: %s", err, out)
	}
	if got := Canonical(junction); got != target {
		t.Errorf("got %q, want %q", got, target)
	}
}
