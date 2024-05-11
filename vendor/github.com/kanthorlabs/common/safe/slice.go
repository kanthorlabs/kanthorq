package safe

import "sync"

type Slice[T any] struct {
	mu   sync.Mutex
	data []T
}

func (sm *Slice[T]) Append(values ...T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.data = append(sm.data, values...)
}

func (sm *Slice[T]) Count() int {
	return len(sm.data)
}

func (sm *Slice[T]) Data() []T {
	return sm.data
}
