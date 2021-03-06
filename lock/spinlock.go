package lock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

const maxBackoff = 64

func (sl *spinLock) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		for i := 0; i < backoff; i++ {
			// yields processor, let other go routines run
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

// NewSpinLock return new lock
func NewSpinLock() sync.Locker {
	return new(spinLock)
}
