package util

import (
	"fmt"
	"os"
	"sync"

	"github.com/sasha-s/go-deadlock"
)

func init() {
	defaultLockType = SyncLock
	deadlock.Opts.DisableLockOrderDetection = false

	switch os.Getenv("KUTIL_LOCK_DEFAULT") {
	case "sync":
		defaultLockType = SyncLock
	case "debug":
		defaultLockType = DebugLock
	case "mock":
		defaultLockType = MockLock
	}

	switch os.Getenv("KUTIL_LOCK_DEBUG_ORDER_DETECTION") {
	case "true":
		deadlock.Opts.DisableLockOrderDetection = false
	case "false":
		deadlock.Opts.DisableLockOrderDetection = true
	}
}

//
// RWLocker
//

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
	RLocker() sync.Locker
}

type LockType int

const (
	DefaultLock LockType = 0
	SyncLock    LockType = 1
	DebugLock   LockType = 2
	MockLock    LockType = 3
)

var defaultLockType LockType

func NewRWLocker(type_ LockType) RWLocker {
	switch type_ {
	case DefaultLock:
		return NewDefaultRWLocker()
	case SyncLock:
		return NewSyncRWLocker()
	case DebugLock:
		return NewDebugRWLocker()
	case MockLock:
		return NewMockRWLocker()
	default:
		panic(fmt.Sprintf("unsupported lock type: %d", type_))
	}
}

//
// DefaultRWLocker
//

func NewDefaultRWLocker() RWLocker {
	switch defaultLockType {
	case SyncLock:
		return NewSyncRWLocker()
	case DebugLock:
		return NewDebugRWLocker()
	case MockLock:
		return NewMockRWLocker()
	default:
		panic(fmt.Sprintf("unsupported lock type: %d", defaultLockType))
	}
}

//
// SyncRWLocker
//

func NewSyncRWLocker() RWLocker {
	return new(sync.RWMutex)
}

//
// DebugRWLocker
//

func NewDebugRWLocker() RWLocker {
	return new(deadlock.RWMutex)
}

//
// MockLocker
//

type MockLocker struct{}

func NewMockLocker() sync.Locker {
	return MockLocker{}
}

// ([sync.Locker] interface)
func (self MockLocker) Lock() {}

// ([sync.Locker] interface)
func (self MockLocker) Unlock() {}

//
// MockRWLocker
//

type MockRWLocker struct {
	MockLocker
}

func NewMockRWLocker() RWLocker {
	return MockRWLocker{}
}

// ([RWLocker] interface)
func (self MockRWLocker) RLock() {}

// ([RWLocker] interface)
func (self MockRWLocker) RUnlock() {}

// ([RWLocker] interface)
func (self MockRWLocker) RLocker() sync.Locker {
	return self
}

//
// LockableEntity
//

type LockableEntity interface {
	GetEntityLock() RWLocker
}

// From [LockableEntity] interface.
func GetEntityLock(entity any) RWLocker {
	if lockable, ok := entity.(LockableEntity); ok {
		return lockable.GetEntityLock()
	} else {
		return nil
	}
}

//
// Ad-hoc locks
//

var adHocLocks sync.Map

// Warning: Because pointers can be re-used after the resource is freed,
// there is no way for us to guarantee ad-hoc locks would not be reused
// Thus this facililty should only be used for objects with a known and managed life span.
func GetAdHocLock(pointer any, type_ LockType) RWLocker {
	if pointer == nil {
		panic("no ad-hoc lock for nil")
	}

	if lock, ok := adHocLocks.Load(pointer); ok {
		return lock.(RWLocker)
	} else {
		lock := NewRWLocker(type_)
		if existing, loaded := adHocLocks.LoadOrStore(pointer, lock); loaded {
			return existing.(RWLocker)
		} else {
			return lock
		}
	}
}

func ResetAdHocLocks() {
	// See: https://stackoverflow.com/a/49355523
	adHocLocks.Range(func(key any, value any) bool {
		adHocLocks.Delete(key)
		return true
	})
}
