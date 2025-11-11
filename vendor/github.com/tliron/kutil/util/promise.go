package util

import (
	contextpkg "context"
	"sync"
)

//
// Promise
//

type Promise chan struct{}

func NewPromise() Promise {
	return make(Promise)
}

func (self Promise) Release() {
	close(self)
}

func (self Promise) Wait(context contextpkg.Context) error {
	select {
	case <-context.Done():
		return context.Err()
	case <-self:
		return nil
	}
}

//
// CoordinatedWork
//

type CoordinatedWork struct {
	sync.Map
}

func NewCoordinatedWork() *CoordinatedWork {
	return new(CoordinatedWork)
}

func (self *CoordinatedWork) Start(context contextpkg.Context, key string) (Promise, bool) {
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		return promise, true
	} else {
		promise = existing.(Promise)
		promise.Wait(context)
		return nil, false
	}
}
