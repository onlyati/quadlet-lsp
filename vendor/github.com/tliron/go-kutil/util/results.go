package util

import (
	"fmt"
	"io"
)

var (
	MaxResultsSliceSize     = 1_000 // 0 for limitless
	ResultsStreamBufferSize = 100
)

type IterateResultsFunc[E any] func(entity E) error
type GetResults[E any] func(offset uint) (Results[E], error)

func IterateResults[E any](results Results[E], iterate IterateResultsFunc[E]) error {
	defer results.Release()

	for {
		if entity, err := results.Next(); err == nil {
			if err := iterate(entity); err != nil {
				return err
			}
		} else if err == io.EOF {
			return nil
		} else {
			return err
		}
	}
}

func GatherResults[E any](results Results[E]) ([]E, error) {
	var slice []E

	if err := IterateResults(results, func(entity E) error {
		slice = append(slice, entity)
		if MaxResultsSliceSize != 0 {
			if len(slice) > MaxResultsSliceSize {
				return fmt.Errorf("too many results: %d", MaxResultsSliceSize)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return slice, nil
}

func CombineResults[E any](get GetResults[E]) Results[E] {
	stream := NewResultsStream[E](nil)

	go func() {
		var offset uint

		for {
			var count uint

			if results, err := get(offset); err == nil {
				if err := IterateResults(results, func(entity E) error {
					stream.Send(entity)
					count++
					return nil
				}); err != nil {
					stream.Close(err)
				}
			} else {
				stream.Close(err)
			}

			if count == 0 {
				break
			}

			offset += count
		}

		stream.Close(nil)
	}()

	return stream
}

//
// Results
//

type Results[E any] interface {
	Next() (E, error) // can return io.EOF
	Release()
}

//
// ResultsSlice
//

type ResultsSlice[E any] struct {
	entities []E
	length   int
	index    int
}

func NewResultsSlice[E any](entities []E) *ResultsSlice[E] {
	return &ResultsSlice[E]{
		entities: entities,
		length:   len(entities),
	}
}

func NewResult[E any](entity E) *ResultsSlice[E] {
	return &ResultsSlice[E]{
		entities: []E{entity},
		length:   1,
	}
}

// ([Results] interface)
func (self *ResultsSlice[E]) Next() (E, error) {
	if self.index < self.length {
		entity := self.entities[self.index]
		self.index++
		return entity, nil
	} else {
		return *new(E), io.EOF
	}
}

// ([Results] interface)
func (self *ResultsSlice[E]) Release() {
}

//
// ResultsStream
//

type ResultsStream[E any] struct {
	release  func()
	entities chan E
	errors   chan error
}

func NewResultsStream[E any](release func()) *ResultsStream[E] {
	return &ResultsStream[E]{
		release:  release,
		entities: make(chan E, ResultsStreamBufferSize),
		errors:   make(chan error),
	}
}

// ([Results] interface)
func (self *ResultsStream[E]) Next() (E, error) {
	for {
		select {
		case entity, ok := <-self.entities:
			if ok {
				return entity, nil
			} else {
				return *new(E), io.EOF
			}

		case err := <-self.errors:
			return *new(E), err
		}
	}
}

// ([Results] interface)
func (self *ResultsStream[E]) Release() {
	if self.release != nil {
		self.release()
	}
}

func (self *ResultsStream[E]) Send(info E) {
	self.entities <- info
}

// Special handling for nil and [io.EOF]
func (self *ResultsStream[E]) Close(err error) {
	if (err == nil) || (err == io.EOF) {
		close(self.entities)
	} else {
		self.errors <- err
	}
}
