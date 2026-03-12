package domain

import "time"

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusLocked    SeatStatus = "LOCKED"
	SeatStatusBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID          string     `json:"id"`
	TrainID     string     `json:"train_id"`
	SeatNumber  string     `json:"seat_number"`
	Status      SeatStatus `json:"status"`
	LockedBy    string     `json:"locked_by_user_id,omitempty"` // Who holds the lock
	LockedUntil *time.Time `json:"locked_until,omitempty"`      // When the lock expires
}
