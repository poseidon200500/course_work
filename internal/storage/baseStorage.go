package storage

import (
	"sync"
)

type BaseStorage struct {
	mu sync.RWMutex

	data          []string
	totalInserted int
}

func NewBaseStorage() Storage {
	return &BaseStorage{
		data: make([]string, 0),
	}
}

// Add — просто добавляем строку без оптимизации
func (s *BaseStorage) Add(str string) {
	s.mu.Lock()
	s.data = append(s.data, str)
	s.totalInserted++
	s.mu.Unlock()
}

func (s *BaseStorage) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.data))
	copy(result, s.data)
	return result
}

func (s *BaseStorage) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return Stats{
		TotalInserted: s.totalInserted,
		UniqueCount:   countUnique(s.data),
	}
}

func (s *BaseStorage) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]string, 0)
	s.totalInserted = 0
}

// вспомогательная функция (дорогая!)
func countUnique(data []string) int {
	set := make(map[string]struct{}, len(data))
	for _, v := range data {
		set[v] = struct{}{}
	}
	return len(set)
}
