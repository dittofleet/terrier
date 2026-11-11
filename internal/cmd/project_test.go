package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dittofleet/terrier/internal/store"
)

// registry builds a store from slash-separated paths, spelled the way this
// OS spells them.
func registry(paths ...string) *store.Store {
	s := &store.Store{}
	for _, p := range paths {
		s.Projects = append(s.Projects, filepath.FromSlash(p))
	}
	return s
}

func TestResolveRefByName(t *testing.T) {
	s := registry("/s/dittofleet/whatagain", "/s/cli/terrier")
	got, err := resolveRef(s, "whatagain")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.FromSlash("/s/dittofleet/whatagain") {
		t.Errorf("got %q", got)
	}
}

func TestResolveRefNeedsAWholeSegment(t *testing.T) {
	s := &store.Store{Projects: []string{"/s/dittofleet/whatagain"}}
	if _, err := resolveRef(s, "again"); err == nil {
		t.Error("a partial segment matched")
	}
}

func TestResolveRefDisambiguatesByPath(t *testing.T) {
	s := registry("/s/alpha/api", "/s/beta/api")
	if _, err := resolveRef(s, "api"); err == nil {
		t.Fatal("an ambiguous ref resolved")
	} else if !strings.Contains(err.Error(), filepath.FromSlash("alpha/api")) || !strings.Contains(err.Error(), filepath.FromSlash("beta/api")) {
		t.Errorf("the error does not name the candidates: %v", err)
	}
	got, err := resolveRef(s, "beta/api")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.FromSlash("/s/beta/api") {
		t.Errorf("got %q", got)
	}
}

func TestResolveRefFallsBackToCaseInsensitive(t *testing.T) {
	s := registry("/s/cli/PReviewer")
	got, err := resolveRef(s, "previewer")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.FromSlash("/s/cli/PReviewer") {
		t.Errorf("got %q", got)
	}
}

func TestResolveRefPrefersAnExactCaseMatch(t *testing.T) {
	s := registry("/s/one/api", "/s/two/API")
	got, err := resolveRef(s, "api")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.FromSlash("/s/one/api") {
		t.Errorf("got %q, want the exactly-cased match", got)
	}
}

func TestResolveRefUnknown(t *testing.T) {
	s := &store.Store{Projects: []string{"/s/cli/terrier"}}
	if _, err := resolveRef(s, "nope"); err == nil {
		t.Error("an unregistered name resolved")
	}
}

func TestLooksLikePath(t *testing.T) {
	cases := map[string]bool{
		"/abs/path": true,
		"~/repos/x": true,
		"./here":    true,
		"../there":  true,
		"whatagain": false,
		"a/b":       false,
	}
	for ref, want := range cases {
		if got := looksLikePath(ref); got != want {
			t.Errorf("looksLikePath(%q) = %v, want %v", ref, got, want)
		}
	}
}

func TestExpandTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	cases := map[string]string{
		"~":       home,
		"~/repos": filepath.Join(home, "repos"),
		"~user":   "~user",
		"repos":   "repos",
	}
	for path, want := range cases {
		if got := expandTilde(path); got != want {
			t.Errorf("expandTilde(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestDescribeReportsAMissingPath(t *testing.T) {
	p := describe("/no/such/path/anywhere")
	if !p.Missing {
		t.Error("a path that does not exist was not reported missing")
	}
	if p.Slug != "" {
		t.Errorf("slug = %q, want empty for a path that is not there", p.Slug)
	}
}

// A repo with no GitHub origin reports a path and nothing else, rather
// than a slug that would mean nothing.
func TestDescribeOmitsSlugWithoutAGitHubOrigin(t *testing.T) {
	dir := t.TempDir()
	p := describe(dir)
	if p.Missing {
		t.Error("an existing directory was reported missing")
	}
	if p.Slug != "" {
		t.Errorf("slug = %q, want empty", p.Slug)
	}
	if p.Path != dir {
		t.Errorf("path = %q, want %q", p.Path, dir)
	}
}
