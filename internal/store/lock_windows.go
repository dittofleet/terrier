package store

import (
	"os"

	"golang.org/x/sys/windows"
)

// errLocked is what tryLock fails with while another process holds the
// lock.
const errLocked = windows.ERROR_LOCK_VIOLATION

// tryLock takes an exclusive lock on f without waiting. Closing f releases
// it.
func tryLock(f *os.File) error {
	const flags = windows.LOCKFILE_EXCLUSIVE_LOCK | windows.LOCKFILE_FAIL_IMMEDIATELY
	return windows.LockFileEx(windows.Handle(f.Fd()), flags, 0, 1, 0, new(windows.Overlapped))
}
