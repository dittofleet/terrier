package store

import "testing"

func TestContainingAtADriveRoot(t *testing.T) {
	s := &Store{Projects: []string{`C:\`}}
	if got := s.Containing(`C:\src\app`); got != `C:\` {
		t.Errorf("got %q, want %q", got, `C:\`)
	}
	if got := s.Containing(`D:\src\app`); got != "" {
		t.Errorf("another drive resolved to %q", got)
	}
}
