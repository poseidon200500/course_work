package storage

import (
	"sync"

	"unique"

	"github.com/poseidon200500/course_work/storage"
)

type UniqueTag struct {
	ptr unique.Handle[string]
	str string
}

type UniqueStorageV2 struct {
	mu sync.RWMutex

	data          []UniqueTag
	totalInserted int
}

func NewUniqueStorageV2() *Storage {
	return &UniqueStorageV2{
		data: make([]UniqueTag, 0),
	}
}

func (s *UniqueStorageV2) Add(str string) {
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

func (s *UniqueStorageV2) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.data))
	for i, t := range s.data {
		result[i] = t.str
	}
	return result
}

func (s *UniqueStorageV2) Stats() Stats {
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

func (s *UniqueStorageV2) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]UniqueTag, 0)
	s.totalInserted = 0
}
