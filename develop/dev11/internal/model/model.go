package model

import (
	"fmt"
	"time"
)

const parseDateFormat string = "2006-01-02"

type (
	Date time.Time

	Event struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		Date        Date   `json:"date"`
		Description string `json:"description"`
	}

	EventCreateRequest struct {
		UserID      int
		Date        Date
		Description string
	}

	EventUpdateRequest struct {
		ID          int
		UserID      int
		Date        Date
		Description string
	}

	EventDeleteRequest struct {
		ID int
	}
	// EventFilter is used to create dates to describe range for requests
	EventFilter struct {
		ID     int
		UserID int
		From   Date
		To     Date
	}
)

func (d Date) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, time.Time(d).Format(parseDateFormat))
	return []byte(s), nil
}

func (ef EventFilter) Match(e Event) bool {
	// If id/user_id is included in filter => proceed to compare individial fields one at a time
	// If id of requested event is not matching filter => data not found
	if ef.ID != 0 {
		if ef.ID != e.ID {
			return false
		}
	}
	// If user_id is not matching => data not found either
	if ef.UserID != 0 {
		if ef.UserID != e.UserID {
			return false
		}
	}
	// If event went before or after specified range => data not found
	f := time.Time(ef.From)
	// If before/after field is included in filter => check
	if !f.IsZero() {
		if f.After(time.Time(e.Date)) {
			return false
		}
	}
	// Same logic for before
	t := time.Time(ef.To)
	if !t.IsZero() {
		if t.Before(time.Time(e.Date)) {
			return false
		}
	}
	return true
}
