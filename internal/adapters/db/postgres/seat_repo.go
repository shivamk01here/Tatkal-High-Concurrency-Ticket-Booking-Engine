package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivamk01here/tatkal-engine/internal/core/domain"
	"github.com/shivamk01here/tatkal-engine/internal/core/ports"
)

type seatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) ports.SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) GetAvailableSeats(ctx context.Context, trainID string, limit int) ([]*domain.Seat, error) {
	query := `
		SELECT id, train_id, seat_number, status 
		FROM seats 
		WHERE train_id = $1 AND status = $2 
		LIMIT $3 
		FOR UPDATE SKIP LOCKED`

	rows, err := r.db.Query(ctx, query, trainID, domain.SeatStatusAvailable, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query available seats: %w", err)
	}
	defer rows.Close()

	var seats []*domain.Seat
	for rows.Next() {
		seat := &domain.Seat{}
		if err := rows.Scan(&seat.ID, &seat.TrainID, &seat.SeatNumber, &seat.Status); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

func (r *seatRepository) UpdateSeatStatus(ctx context.Context, seatID string, status domain.SeatStatus) error {
	query := `UPDATE seats SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, seatID)
	if err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}
	return nil
}
