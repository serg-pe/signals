package array

import (
	"errors"
	"sync/atomic"

	"github.com/serg-pe/signals/pkg/types/queue"
)

var (
	errNegativeRangeNotSupported = errors.New("negative range is not supported")
	errOutOfRange                = errors.New("out of range")
	errAlreadyRemoved            = errors.New("already removed")
)

type stored[T any] struct {
	used bool
	data T
}

type ArrayStorage[T any] struct {
	// TODO: mutex
	// mu     *sync.RWMutex

	storage  []stored[T]
	length   int
	capacity int
	freed    queue.Queue[int]
	unuseds  int
}

func New[T any](capacity int) ArrayStorage[T] {
	data := make([]stored[T], capacity)

	lock := &atomic.Bool{}
	lock.Store(false)

	return ArrayStorage[T]{
		storage:  data,
		length:   0,
		capacity: cap(data),
		freed:    queue.New[int](),
		unuseds:  0,
	}
}

func (s *ArrayStorage[T]) Add(entry T) int {
	var id int

	if id, ok := s.freed.Pop(); ok {
		s.unuseds--
		s.storage[id].used = true
		s.storage[id].data = entry
		return id
	}

	if s.length < s.capacity {
		id = s.length
		s.storage[id] = stored[T]{true, entry}
		s.length++
		return id
	}

	return s.addWithReallocation(entry)
}

func (s *ArrayStorage[T]) Get(id int) (T, error) {
	var result T

	if err := s.checkBounds(id); err != nil {
		return result, err
	}

	if !s.storage[id].used {
		return result, errAlreadyRemoved
	}

	result = s.storage[id].data
	return result, nil
}

func (s *ArrayStorage[T]) Remove(id int) error {
	if err := s.checkBounds(id); err != nil {
		return err
	}

	if !s.storage[id].used {
		return errAlreadyRemoved
	}

	s.storage[id].used = false
	if id == s.length-1 {
		s.length--
		return nil
	}

	s.freed.Push(id)
	s.unuseds++
	return nil
}

func (s *ArrayStorage[T]) ApplyTo(filter func(entry T) bool, apply func(entry T) T) {
	for index := 0; index < s.length; index++ {
		if s.storage[index].used && filter(s.storage[index].data) {
			s.storage[index].data = apply(s.storage[index].data)
		}
	}
}

func (s *ArrayStorage[T]) ApplyToAll(apply func(entry T)) {
	for index := 0; index < s.length; index++ {
		if s.storage[index].used {
			apply(s.storage[index].data)
		}
	}
}

func (s *ArrayStorage[T]) Update(id int, update func(entry T) T) error {
	entry, err := s.Get(id)
	if err != nil {
		return err
	}

	s.storage[id].data = update(entry)
	return nil
}

func (s *ArrayStorage[T]) addWithReallocation(entry T) int {
	id := s.length
	s.storage = append(s.storage, stored[T]{true, entry})
	s.capacity = cap(s.storage)
	s.length++
	return id
}

func (s *ArrayStorage[T]) checkBounds(id int) error {
	if id < 0 {
		return errNegativeRangeNotSupported
	}

	if id > s.length-1 {
		return errOutOfRange
	}

	return nil
}
