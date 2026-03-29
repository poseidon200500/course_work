package uniquestr

import (
	"sync"

	"your_project/storage"

	"unique"
)

type UniqueTag struct {
	ptr unique.Handle[string]
	str string
}

type Storage struct {
	mu sync.RWMutex

	data          []UniqueTag
	totalInserted int
}

func NewStorage() *Storage {
	return &Storage{
		data: make([]UniqueTag, 0),
	}
}

func (s *Storage) Add(str string) {
	h := unique.Make(str)

	tag := UniqueTag{
		ptr: h,
		str: h.Value(),
	}

	s.mu.Lock()
	s.data = append(s.data, tag)
	s.totalInserted++
	s.mu.Unlock()
}

func (s *Storage) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.data))
	for i, t := range s.data {
		result[i] = t.str
	}
	return result
}

func (s *Storage) Stats() storage.Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set := make(map[string]struct{})
	for _, t := range s.data {
		set[t.str] = struct{}{}
	}

	return storage.Stats{
		TotalInserted: s.totalInserted,
		UniqueCount:   len(set),
	}
}

func (s *Storage) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]UniqueTag, 0)
	s.totalInserted = 0
}
