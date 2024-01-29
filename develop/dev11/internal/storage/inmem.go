package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/Karanth1r3/l_2/develop/dev11/internal/model"
)

var ErrNotFound = errors.New("not found")

// Data is saved to map in inmem storage variant
type (
	InMemStorage struct {
		sid  int
		data map[int]model.Event
		mu   sync.RWMutex
	}
)

// Constructor
func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		sid:  1,
		data: make(map[int]model.Event),
		mu:   sync.RWMutex{},
	}
}

// (s *InMemStorage) CreateEvent(...) handles id initialization logic & new event creation
func (s *InMemStorage) CreateEvent(_ context.Context, r model.EventCreateRequest) (int, error) {

	s.mu.RLock()
	id := s.sid
	s.data[id] = model.Event{
		ID:          id,
		UserID:      r.UserID,
		Date:        r.Date,
		Description: r.Description,
	}
	s.sid++
	s.mu.RUnlock()
	return id, nil
}

func (s *InMemStorage) UpdateEvent(_ context.Context, r model.EventUpdateRequest) error {
	s.mu.RLock()
	// If there is no record with id specified in request => return not found
	if _, ok := s.data[r.ID]; !ok {
		return ErrNotFound
	}
	// Update event data with request fields
	s.data[r.ID] = model.Event(r)
	s.mu.RUnlock()
	return nil
}

func (s *InMemStorage) DeleteEvent(_ context.Context, r model.EventDeleteRequest) error {
	s.mu.RLock()
	// If there is no record with id specified in request => return not found
	if _, ok := s.data[r.ID]; !ok {
		return ErrNotFound
	}
	// Delete event by id
	delete(s.data, r.ID)
	// To prevent overwrite on existing record (that's quite possible actually)
	s.data[-r.ID] = model.Event{}
	s.mu.RUnlock()

	return nil
}

func (s *InMemStorage) GetEvents(_ context.Context, ef model.EventFilter) ([]model.Event, error) {

	result := make([]model.Event, 0)
	// If event from storage is within provided datetime range => add it to result
	for _, v := range s.data {
		if ef.Match(v) {
			result = append(result, v)
		}
	}
	// TODO: Probably should handle len 0; Not sure yeat
	return result, nil
}
