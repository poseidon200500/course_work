package storage

import (
	"sync"
	"unique"
)

type Storage struct {
	mu sync.RWMutex

	data          []unique.Handle[string]
	totalInserted int
}

func NewStorage() *Storage {
	return &Storage{
		data: make([]unique.Handle[string], 0),
	}
}

func (s *Storage) Add(str string) {
	h := unique.Make(str)

	s.mu.Lock()
	s.data = append(s.data, h)
	s.totalInserted++
	s.mu.Unlock()
}

func (s *Storage) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.data))
	for i, h := range s.data {
		result[i] = h.Value()
	}
	return result
}

func (s *Storage) Stats() Stats {
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

func (s *Storage) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]unique.Handle[string], 0)
	s.totalInserted = 0
}
