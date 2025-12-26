package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/arkine/l4/3/internal/calendar"
)

type Handler struct {
	Service *calendar.Service
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int    `json:"user_id"`
		Date     string `json:"date"`
		Text     string `json:"text"`
		RemindAt string `json:"remind_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format"}`, http.StatusBadRequest)
		return
	}

	var remind *time.Time
	if req.RemindAt != "" {
		t, err := time.Parse(time.RFC3339, req.RemindAt)
		if err != nil {
			http.Error(w, `{"error": "invalid remind_at format"}`, http.StatusBadRequest)
			return
		}
		remind = &t
	}

	ev, err := h.Service.Create(calendar.Event{
		UserID:   req.UserID,
		Date:     date,
		Text:     req.Text,
		RemindAt: remind,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": ev})
	if err != nil {
		fmt.Printf(`{"error": "%s"}`, err)
	}
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID, _ := strconv.Atoi(q.Get("user_id"))
	dateStr := q.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
		return
	}

	events := h.Service.EventsForDay(userID, date)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if err != nil {
		fmt.Printf(`{"error": "%s"}`, err)
	}
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID, _ := strconv.Atoi(q.Get("user_id"))
	dateStr := q.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
		return
	}

	events := h.Service.EventsForWeek(userID, date)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if err != nil {
		fmt.Printf(`{"error": "%s"}`, err)
	}
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID, _ := strconv.Atoi(q.Get("user_id"))
	dateStr := q.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
		return
	}

	events := h.Service.EventsForMonth(userID, date)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if err != nil {
		fmt.Printf("Error encoding events: %v\n", err)
	}
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       int    `json:"id"`
		UserID   int    `json:"user_id"`
		Date     string `json:"date"`
		Text     string `json:"text"`
		RemindAt string `json:"remind_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format"}`, http.StatusBadRequest)
		return
	}

	var remind *time.Time
	if req.RemindAt != "" {
		t, err := time.Parse(time.RFC3339, req.RemindAt)
		if err != nil {
			http.Error(w, `{"error": "invalid remind_at format"}`, http.StatusBadRequest)
			return
		}
		remind = &t
	}

	err = h.Service.Update(req.ID, calendar.Event{
		ID:       req.ID,
		UserID:   req.UserID,
		Date:     date,
		Text:     req.Text,
		RemindAt: remind,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	if err != nil {
		fmt.Printf("Error encoding events: %v\n", err)
	}
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	err := h.Service.Delete(req.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	if err != nil {
		fmt.Printf("Error encoding events: %v\n", err)
	}
}
