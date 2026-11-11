package store

import (
	"errors"
	"os"
	"testing"
)

func TestLockExcludesASecondHolderUntilReleased(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	release, err := acquire()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(lockPath(), os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := tryLock(f); !errors.Is(err, errLocked) {
		t.Fatalf("locking a held lock: got %v, want %v", err, errLocked)
	}

	release()
	if err := tryLock(f); err != nil {
		t.Errorf("locking after release: %v", err)
	}
}
