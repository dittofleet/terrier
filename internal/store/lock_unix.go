//go:build !windows

package store

import (
	"os"
	"syscall"
)

// errLocked is what tryLock fails with while another process holds the
// lock.
const errLocked = syscall.EWOULDBLOCK

// tryLock takes an exclusive lock on f without waiting. Closing f releases
// it.
func tryLock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}
