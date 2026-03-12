package domain

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusFailed    BookingStatus = "FAILED"
)

type Booking struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	TrainID   string        `json:"train_id"`
	SeatIDs   []string      `json:"seat_ids"`
	Status    BookingStatus `json:"status"`
	Amount    float64       `json:"amount"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
