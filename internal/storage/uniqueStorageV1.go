package storage

import (
	"sync"
	"unique"
)

type UniqueStorageV1 struct {
	mu sync.RWMutex

	data          []unique.Handle[string]
	totalInserted int
}

func NewUniqueStorageV1() Storage {
	return &UniqueStorageV1{
		data: make([]unique.Handle[string], 0),
	}
}

func (s *UniqueStorageV1) Add(str string) {
	h := unique.Make(str)

	s.mu.Lock()
	s.data = append(s.data, h)
	s.totalInserted++
	s.mu.Unlock()
}

func (s *UniqueStorageV1) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.data))
	for i, h := range s.data {
		result[i] = h.Value()
	}
	return result
}

func (s *UniqueStorageV1) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// уникальные считаем через map (дорого, но честно)
	set := make(map[string]struct{})
	for _, h := range s.data {
		set[h.Value()] = struct{}{}
	}

	return Stats{
		TotalInserted: s.totalInserted,
		UniqueCount:   len(set),
	}
}

func (s *UniqueStorageV1) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]unique.Handle[string], 0)
	s.totalInserted = 0
}
