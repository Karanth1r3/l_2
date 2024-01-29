package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Karanth1r3/l_2/develop/dev11/internal/model"
)

type (
	Handler struct {
		service calendarService
	}

	calendarService interface {
		CreateEvent(ctx context.Context, r model.EventCreateRequest) (model.Event, error)
		UpdateEvent(ctx context.Context, r model.EventUpdateRequest) error
		DeleteEvent(ctx context.Context, r model.EventDeleteRequest) error
		GetEvents(ctx context.Context, r model.EventFilter) ([]model.Event, error)
	}

	resultResp struct {
		Result any `json:"result"`
	}

	errorResp struct {
		Error any `json:"error"`
	}
)

// Ctor
func NewHandler(cs calendarService) *Handler {
	return &Handler{
		service: cs,
	}
}

// Main API handlers for endpoints
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowedError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}

	reqData, err := parseEventCreateRequest(r.Form)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	event, err := h.service.CreateEvent(r.Context(), reqData)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowedError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}

	reqData, err := parseEventUpdateRequest(r.Form)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	err = h.service.UpdateEvent(r.Context(), reqData)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, "Update Success")
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowedError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}

	reqData, err := parseEventDeleteRequest(r.Form)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	err = h.service.DeleteEvent(r.Context(), reqData)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, "Delete Success")
}

// Creates range for request with day filter
func findDayRange(e time.Time) (low, high time.Time) {
	low = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	high = time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
	return low, high
}

// Try to get set of records that correspond to request (in this case => events on day)
func (h *Handler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	// If request has wrong method => Write error & return
	if r.Method != http.MethodGet {
		writeMethodNotAllowedError(w)
		return
	}
	// If request could not be parsed => Write error & return
	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}
	// If data could not be parsed from req params => Write err & return
	userId, date, err := parseDailyEventsRequest(r.Form)
	if err != nil {
		writeBadRequestError(w, err)
		return
	}

	low, high := findDayRange(time.Time(date))
	// Try to get records from storage. If failed => Write about internal server error & return
	result, err := h.service.GetEvents(r.Context(), model.EventFilter{
		UserID: userId,
		From:   model.Date(low),
		To:     model.Date(high),
	})
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, result)
}

// Data to convert weeks starting with Sunday (internal) to weeks starting with Monday (for application)
type tempDay int

const (
	Monday tempDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// Monday will be the first day of the week after convert
func convertWeekDay(t time.Weekday) tempDay {
	val := int(t)
	val--
	if val < 0 {
		val = 7
	}
	converted := tempDay(val)
	return converted
}

// finds the first and the last day of the week
func findWeekRange(e time.Time) (lower, higher time.Time) {
	wDay := (e.Weekday())
	r := convertWeekDay(wDay) // Offset for week to begin with Monday for this case

	var inc, dec int = 0, 0
	for i := r; i > 1; i-- {
		dec++
	}
	for i := r; i < 7; i++ {
		inc++
	}

	lower = e.Add(time.Hour * time.Duration(-24*dec))
	higher = e.Add(time.Hour * time.Duration(24*inc))
	lower = time.Date(lower.Year(), lower.Month(), lower.Day(), 0, 0, 0, 0, e.Location())
	higher = time.Date(higher.Year(), higher.Month(), higher.Day(), 23, 59, 59, 0, e.Location())
	return lower, higher
}

// Try to get set of records that correspond to request (in this case => events on day)
func (h *Handler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	// If request has wrong method => Write error & return
	if r.Method != http.MethodGet {
		writeMethodNotAllowedError(w)
		return
	}
	// If request could not be parsed => Write error & return
	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}
	// If data could not be parsed from req params => Write err & return
	userId, date, err := parseDailyEventsRequest(r.Form)
	if err != nil {
		writeBadRequestError(w, err)
		return
	}

	low, high := findWeekRange(time.Time(date))
	// Try to get records from storage. If failed => Write about internal server error & return
	result, err := h.service.GetEvents(r.Context(), model.EventFilter{
		UserID: userId,
		From:   model.Date(low),
		To:     model.Date(high),
	})
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, result)
}

// daysIn returns number of days in specified month
func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// Gets range for request with month filter to find events
func findMonthRange(e time.Time) (low, high time.Time) {
	month := e.Month()
	inc, dec := 0, 0
	for i := e.Day(); i > 1; i-- {
		dec++
	}

	for i := int(e.Day()); i < daysIn(month, e.Year()); i++ {
		inc++
	}
	low = e.Add(time.Hour * time.Duration(-24*dec))
	high = e.Add(time.Hour * time.Duration(24*inc))
	low = time.Date(low.Year(), low.Month(), low.Day(), 0, 0, 0, 0, e.Location())
	high = time.Date(high.Year(), high.Month(), high.Day(), 23, 59, 59, 0, e.Location())
	return low, high
}

// Try to get set of records that correspond to request (in this case => events on day)
func (h *Handler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	// If request has wrong method => Write error & return
	if r.Method != http.MethodGet {
		writeMethodNotAllowedError(w)
		return
	}
	// If request could not be parsed => Write error & return
	if err := r.ParseForm(); err != nil {
		writeBadRequestError(w, err)
		return
	}
	// If data could not be parsed from req params => Write err & return
	userId, date, err := parseDailyEventsRequest(r.Form)
	if err != nil {
		writeBadRequestError(w, err)
		return
	}

	low, high := findMonthRange(time.Time(date))
	// Try to get records from storage. If failed => Write about internal server error & return
	result, err := h.service.GetEvents(r.Context(), model.EventFilter{
		UserID: userId,
		From:   model.Date(low),
		To:     model.Date(high),
	})
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResult(w, result)
}

// Output utility for response writer SECTION:
func writeResult(w http.ResponseWriter, result any) {
	data, _ := json.Marshal(&resultResp{
		Result: result,
	})
	_, _ = w.Write(data)
}

// Errors form & write section (in the format required by task)
func writeMethodNotAllowedError(w http.ResponseWriter) {
	data, _ := json.Marshal(&errorResp{
		Error: http.StatusMethodNotAllowed,
	})
	writeError(w, data, http.StatusMethodNotAllowed)
}

func writeBadRequestError(w http.ResponseWriter, err error) {
	data, _ := json.Marshal(&errorResp{
		Error: err.Error(),
	})
	writeError(w, data, http.StatusBadRequest)
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	data, _ := json.Marshal(&errorResp{
		Error: err.Error(),
	})
	writeError(w, data, http.StatusInternalServerError)
}

// Write error to request
func writeError(w http.ResponseWriter, body []byte, statusCode int) {
	http.Error(w, string(body), statusCode)
}
