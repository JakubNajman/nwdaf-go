package main

import (
	"fmt"
	"sync"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemStore struct {
	mu     sync.Mutex
	items  map[string]Item
	nextID int
}

func NewItemStore() *ItemStore {
	s := ItemStore{}
	s.items = make(map[string]Item)
	return &s
}

func (s *ItemStore) Add(name string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++

	item := Item{
		ID:   fmt.Sprintf("item-%d", s.nextID),
		Name: name,
	}
	s.items[item.ID] = item
	return item
}

func (s *ItemStore) Get(id string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	return item, ok
}

func (s *ItemStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[id]

	if !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *ItemStore) List() []Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Item, 0, len(s.items))

	for _, item := range s.items {
		out = append(out, item)
	}

	return out
}
