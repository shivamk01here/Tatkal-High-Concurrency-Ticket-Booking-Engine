package ports

import (
	"context"

	"github.com/shivamk01here/tatkal-engine/internal/core/domain"
)

type TrainRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Train, error)
	UpdateAvailableSeats(ctx context.Context, trainID string, seatsToDeduct int) error
}

type SeatRepository interface {
	GetAvailableSeats(ctx context.Context, trainID string, limit int) ([]*domain.Seat, error)
	UpdateSeatStatus(ctx context.Context, seatID string, status domain.SeatStatus) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	UpdateStatus(ctx context.Context, bookingID string, status domain.BookingStatus) error
}
