package storage

import (
	"context"
	"errors"

	"github.com/Karanth1r3/l_2/develop/dev11/internal/model"
)

var ErrNotFound = errors.New("not found")

// Data is saved to map in inmem storage variant
type (
	InMemStorage struct {
		data map[int]model.Event
	}
)

// Constructor
func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		data: make(map[int]model.Event),
	}
}

// (s *InMemStorage) CreateEvent(...) handles id initialization logic & new event creation
func (s *InMemStorage) CreateEvent(_ context.Context, r model.EventCreateRequest) (int, error) {

	id := len(s.data) + 1
	s.data[id] = model.Event{
		ID:          id,
		UserID:      r.UserID,
		Date:        r.Date,
		Description: r.Description,
	}

	return id, nil
}

func (s *InMemStorage) UpdateEvent(_ context.Context, r model.EventUpdateRequest) error {
	// If there is no record with id specified in request => return not found
	if _, ok := s.data[r.ID]; !ok {
		return ErrNotFound
	}
	// Update event data with request fields
	s.data[r.ID] = model.Event(r)

	return nil
}

func (s *InMemStorage) DeleteEvent(_ context.Context, r model.EventDeleteRequest) error {
	// If there is no record with id specified in request => return not found
	if _, ok := s.data[r.ID]; !ok {
		return ErrNotFound
	}
	// Delete event by id
	delete(s.data, r.ID)
	// To prevent overwrite on existing record (that's quite possible actually)
	s.data[-r.ID] = model.Event{}

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
