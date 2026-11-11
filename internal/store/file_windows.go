package store

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// busyTimeout covers another process having the registry open for a moment.
// It is shorter than lockTimeout because a file that is denied for good
// costs every `terrier ls` this long.
const busyTimeout = time.Second

// readFile reads path. Windows refuses to open a file while it is being
// replaced, so a save in progress is waited out.
func readFile(path string) ([]byte, error) {
	var data []byte
	err := retryBusy(func() (err error) {
		data, err = os.ReadFile(path)
		return err
	})
	return data, err
}

// replace puts tmp in path's place. Windows refuses to replace a file
// while it is open, and readers take no lock, so a `terrier ls` that
// happens to be reading is waited out.
func replace(tmp, path string) error {
	return retryBusy(func() error { return os.Rename(tmp, path) })
}

func retryBusy(op func() error) error {
	deadline := time.Now().Add(busyTimeout)
	for {
		err := op()
		busy := errors.Is(err, windows.ERROR_ACCESS_DENIED) || errors.Is(err, windows.ERROR_SHARING_VIOLATION)
		if !busy || time.Now().After(deadline) {
			return err
		}
		time.Sleep(lockPollWait)
	}
}
