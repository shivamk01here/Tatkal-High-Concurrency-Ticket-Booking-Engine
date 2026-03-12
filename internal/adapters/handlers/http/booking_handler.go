package http

import (
	"encoding/json"
	"net/http"

	"github.com/shivamk01here/tatkal-engine/internal/core/services"
)

type BookingHandler struct {
	service *services.TatkalService
}

func NewBookingHandler(service *services.TatkalService) *BookingHandler {
	return &BookingHandler{service: service}
}

type BookingRequest struct {
	UserID    string `json:"user_id"`
	TrainID   string `json:"train_id"`
	SeatCount int    `json:"seat_count"`
}

func (h *BookingHandler) HandleInitiateBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	booking, err := h.service.InitiateBooking(r.Context(), req.UserID, req.TrainID, req.SeatCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict is perfect for Tatkal failures
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}
