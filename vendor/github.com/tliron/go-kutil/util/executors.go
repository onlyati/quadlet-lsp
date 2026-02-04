package util

import (
	"fmt"
	"sync"
)

//
// Executor
//

type Executor[T any] interface {
	Queue(task T)
	Wait() []error
	Close()
}

//
// ParallelExecutor
//

type ParallelExecutor[T any] struct {
	PanicAsError string // when non-empty, wil capture panics as errors under this name

	processor func(task T) error
	tasks     chan T
	wg        sync.WaitGroup
	errs      []error
	errsLock  sync.Mutex
}

func NewParallelExecutor[T any](bufferSize int, processor func(task T) error) *ParallelExecutor[T] {
	self := ParallelExecutor[T]{
		processor: processor,
		tasks:     make(chan T, bufferSize),
	}

	return &self
}

func (self *ParallelExecutor[T]) Start(workers int) {
	for range workers {
		go self.worker()
	}
}

// ([Executor] interface)
func (self *ParallelExecutor[T]) Queue(task T) {
	self.wg.Add(1)
	self.tasks <- task
}

// ([Executor] interface)
func (self *ParallelExecutor[T]) Wait() []error {
	self.wg.Wait()
	close(self.tasks)
	return self.errs
}

// ([Executor] interface)
func (self *ParallelExecutor[T]) Close() {
	close(self.tasks)
}

func (self *ParallelExecutor[T]) worker() {
	for {
		select {
		case task, ok := <-self.tasks:
			if ok {
				if err := self.process(task); err != nil {
					self.errsLock.Lock()
					self.errs = append(self.errs, err)
					self.errsLock.Unlock()
				}
				self.wg.Done()
			} else {
				break
			}
		}
	}
}

func (self *ParallelExecutor[T]) process(task T) (rerr error) {
	if self.PanicAsError != "" {
		defer func() {
			if err := recover(); err != nil {
				rerr = fmt.Errorf("panic during %s: %v", self.PanicAsError, err)
			}
		}()
	}

	return self.processor(task)
}
