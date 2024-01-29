package service

import (
	"context"
	"fmt"

	"github.com/Karanth1r3/l_2/develop/dev11/internal/model"
)

type (
	Service struct {
		storage eventStorage
	}

	eventStorage interface {
		// Return event id or error for creation event, error is enough for other cases
		CreateEvent(context.Context, model.EventCreateRequest) (int, error)
		UpdateEvent(context.Context, model.EventUpdateRequest) error
		DeleteEvent(context.Context, model.EventDeleteRequest) error
		// Filter is specified by request (events for day or month - for example).
		// If there are any events saved within range specified by filter => return them or error
		GetEvents(context.Context, model.EventFilter) ([]model.Event, error)
	}
)

func New(es eventStorage) *Service {
	return &Service{storage: es}
}

func (s *Service) CreateEvent(ctx context.Context, r model.EventCreateRequest) (model.Event, error) {
	// If user_id in request is not initialized => request is not valid
	if r.UserID == 0 {
		return model.Event{}, fmt.Errorf("user_id is not valid")
	}

	id, err := s.storage.CreateEvent(ctx, r)
	if err != nil {
		return model.Event{}, err
	}

	return model.Event{
		ID:          id,
		UserID:      r.UserID,
		Date:        r.Date,
		Description: r.Description,
	}, nil
}

// (s *Service) UpdateEvent() Calls internal storage UpdateEvent() method.
// If data within selected event was successfully updated => nil will be returned instead of error
func (s *Service) UpdateEvent(ctx context.Context, r model.EventUpdateRequest) error {
	return s.storage.UpdateEvent(ctx, r)
}

// (s *Service) DeleteEvent() calls linked storage DeleteEvent() method to handle selected event delete
func (s *Service) DeleteEvent(ctx context.Context, r model.EventDeleteRequest) error {
	return s.storage.DeleteEvent(ctx, r)
}

// Returns events based on specified date filter (month, day, week)
// Calls storage linked to the service to handle request
func (s *Service) GetEvents(ctx context.Context, ef model.EventFilter) ([]model.Event, error) {
	return s.storage.GetEvents(ctx, ef)
}
