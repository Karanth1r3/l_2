package httpapi

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Karanth1r3/l_2/develop/dev11/internal/model"
)

const (
	queryParamEventID     = "id"
	queryParamUserID      = "user_id"
	queryParamDate        = "date"
	queryParamDescription = "description"

	queryDateFormat = "2006-01-02"
)

type (
	eventParams struct {
		UserID      int
		Date        model.Date
		Description string
	}
)

func parseEventCreateRequest(vals url.Values) (empty model.EventCreateRequest, err error) {
	ep, err := parseEventParams(vals)
	if err != nil {
		return empty, nil
	}

	return model.EventCreateRequest{
		UserID:      ep.UserID,
		Date:        ep.Date,
		Description: ep.Description,
	}, nil
}

func parseEventUpdateRequest(vals url.Values) (empty model.EventUpdateRequest, err error) {

	id, err := parseEventID(vals)
	if err != nil {
		return empty, err
	}

	ep, err := parseEventParams(vals)
	if err != nil {
		return empty, nil
	}

	return model.EventUpdateRequest{
		ID:          id,
		UserID:      ep.UserID,
		Date:        ep.Date,
		Description: ep.Description,
	}, nil
}

func parseEventDeleteRequest(vals url.Values) (empty model.EventDeleteRequest, err error) {

	id, err := parseEventID(vals)
	if err != nil {
		return empty, nil
	}

	return model.EventDeleteRequest{
		ID: id,
	}, nil
}

func parseDailyEventsRequest(vals url.Values) (userID int, date model.Date, err error) {
	return parseCommonGetEventsRequest(vals)
}

func parseMonthlyEventsRequest(vals url.Values) (userID int, date model.Date, err error) {
	return parseCommonGetEventsRequest(vals)
}

func parseWeeklyEventsRequest(vals url.Values) (userID int, date model.Date, err error) {
	return parseCommonGetEventsRequest(vals)
}

func parseCommonGetEventsRequest(vals url.Values) (userID int, date model.Date, err error) {
	var ep eventParams
	ep, err = parseEventParams(vals)
	if err != nil {
		return 0, date, err
	}

	return ep.UserID, ep.Date, nil
}

func parseUserID(vals url.Values) (int, error) {
	// Trying to parse userID from url params
	userID, err := strconv.Atoi(vals.Get(queryParamUserID))
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format: %w", err)
	}

	return userID, nil
}

func parseEventID(vals url.Values) (int, error) {
	// Trying to parse eventID from url params
	eventID, err := strconv.Atoi(vals.Get(queryParamEventID))
	if err != nil {
		return 0, fmt.Errorf("invalid event_id format: %w", err)
	}

	return eventID, nil
}

func parseDate(vals url.Values) (model.Date, error) {
	// Trying to parse date from query within format specified in constants section
	date, err := time.Parse(queryDateFormat, queryDateFormat)
	if err != nil {
		return model.Date{}, fmt.Errorf("invalid date format: %w", err)
	}

	return model.Date(date), nil
}

func parseEventParams(vals url.Values) (empty eventParams, err error) {
	// Trying to get data for event objects from urls
	userID, err := parseUserID(vals)
	if err != nil {
		return empty, err
	}

	date, err := parseDate(vals)
	if err != nil {
		return empty, err
	}

	return eventParams{
		UserID:      userID,
		Date:        date,
		Description: vals.Get(queryParamDescription),
	}, nil
}
