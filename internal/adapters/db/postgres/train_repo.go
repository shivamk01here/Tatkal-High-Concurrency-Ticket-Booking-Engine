package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivamk01here/tatkal-engine/internal/core/domain"
	"github.com/shivamk01here/tatkal-engine/internal/core/ports"
)

// trainRepository implements ports.TrainRepository
type trainRepository struct {
	db *pgxpool.Pool
}

// NewTrainRepository creates a new instance of the repository
func NewTrainRepository(db *pgxpool.Pool) ports.TrainRepository {
	return &trainRepository{db: db}
}

func (r *trainRepository) GetByID(ctx context.Context, id string) (*domain.Train, error) {
	query := `
		SELECT id, train_number, name, source, destination, total_seats, available_seats, departure_time 
		FROM trains 
		WHERE id = $1`

	train := &domain.Train{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&train.ID, &train.TrainNumber, &train.Name, &train.Source,
		&train.Destination, &train.TotalSeats, &train.AvailableSeats, &train.DepartureTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get train: %w", err)
	}

	return train, nil
}

func (r *trainRepository) UpdateAvailableSeats(ctx context.Context, trainID string, seatsToDeduct int) error {
	// Notice the highly optimized SQL here. We don't fetch and then update (which causes race conditions).
	// We update it atomically directly in the database.
	query := `
		UPDATE trains 
		SET available_seats = available_seats - $1 
		WHERE id = $2 AND available_seats >= $1`

	tag, err := r.db.Exec(ctx, query, seatsToDeduct, trainID)
	if err != nil {
		return fmt.Errorf("failed to update available seats: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not enough seats available or train not found")
	}

	return nil
}
