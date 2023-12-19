package localstore

import (
	"sync"
)

type Store struct {
	store map[int64][]interface{}

	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		store: make(map[int64][]interface{}, 10),
	}
}

func (s *Store) Set(data []interface{}, userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = data
}

func (s *Store) Read(userID int64) ([]interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.store[userID]
	if !ok {
		return nil, false
	}

	return d, true
}

func (s *Store) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
}
