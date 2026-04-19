package storage

import "sync"

type InternStorage struct {
	mu sync.RWMutex

	pool map[string]string // интернированные строки
	data []string          // все добавленные строки

	totalInserted int
}

func NewInternStorage() Storage {
	return &InternStorage{
		pool: make(map[string]string),
		data: make([]string, 0),
	}
}

// Add с интернированием
func (s *InternStorage) Add(str string) {
	// сначала ищем без блокировки записи
	s.mu.RLock()
	interned, ok := s.pool[str]
	s.mu.RUnlock()

	if !ok {
		s.mu.Lock()

		// double-check
		if existing, exists := s.pool[str]; exists {
			interned = existing
		} else {
			s.pool[str] = str
			interned = str
		}

		s.mu.Unlock()
	}

	// добавляем в data
	s.mu.Lock()
	s.data = append(s.data, interned)
	s.totalInserted++
	s.mu.Unlock()
}

func (s *InternStorage) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// копия (чтобы не сломали извне)
	result := make([]string, len(s.data))
	copy(result, s.data)
	return result
}

func (s *InternStorage) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return Stats{
		TotalInserted: s.totalInserted,
		UniqueCount:   len(s.pool),
	}
}

func (s *InternStorage) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pool = make(map[string]string)
	s.data = make([]string, 0)
	s.totalInserted = 0
}
