package cmd

import (
	"path/filepath"
	"testing"
)

func TestLooksLikeAWindowsPath(t *testing.T) {
	for _, ref := range []string{`C:\repos\x`, `C:/repos/x`, `\repos\x`, `\\server\share\x`, `.\here`, `~\repos\x`} {
		if !looksLikePath(ref) {
			t.Errorf("looksLikePath(%q) = false", ref)
		}
	}
}

func TestExpandTildeTakesABackslash(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	if got, want := expandTilde(`~\repos`), filepath.Join(home, "repos"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
