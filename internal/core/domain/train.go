package domain

import "time"

type Train struct {
	ID             string    `json:"id"`
	TrainNumber    string    `json:"train_number"`
	Name           string    `json:"name"`
	Source         string    `json:"source"`
	Destination    string    `json:"destination"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	DepartureTime  time.Time `json:"departure_time"`
}
