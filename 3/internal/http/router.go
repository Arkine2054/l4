package httpapi

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/create_event", h.CreateEvent)
	router.HandleFunc("/update_event", h.UpdateEvent)
	router.HandleFunc("/delete_event", h.DeleteEvent)
	router.HandleFunc("/events_for_day", h.EventsForDay)
	router.HandleFunc("/events_for_week", h.EventsForWeek)
	router.HandleFunc("/events_for_month", h.EventsForMonth)
	return router
}
