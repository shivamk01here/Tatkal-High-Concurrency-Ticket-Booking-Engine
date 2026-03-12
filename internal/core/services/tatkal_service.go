package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivamk01here/tatkal-engine/internal/core/domain"
	"github.com/shivamk01here/tatkal-engine/internal/core/ports"
)

// TatkalService orchestrates the booking flow
type TatkalService struct {
	trainRepo   ports.TrainRepository
	seatRepo    ports.SeatRepository
	bookingRepo ports.BookingRepository
	cacheRepo   ports.CacheRepository
}

// NewTatkalService injects the dependencies
func NewTatkalService(tr ports.TrainRepository, sr ports.SeatRepository, br ports.BookingRepository, cr ports.CacheRepository) *TatkalService {
	return &TatkalService{
		trainRepo:   tr,
		seatRepo:    sr,
		bookingRepo: br,
		cacheRepo:   cr,
	}
}

// InitiateBooking is the core Tatkal algorithm
func (s *TatkalService) InitiateBooking(ctx context.Context, userID string, trainID string, seatCount int) (*domain.Booking, error) {
	// 1. Check if the train even has enough seats overall (Fast fail)
	train, err := s.trainRepo.GetByID(ctx, trainID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch train details: %w", err)
	}
	if train.AvailableSeats < seatCount {
		return nil, fmt.Errorf("not enough seats available on this train")
	}

	// 2. Fetch specific available seats using our optimized SKIP LOCKED query
	availableSeats, err := s.seatRepo.GetAvailableSeats(ctx, trainID, seatCount)
	if err != nil {
		return nil, fmt.Errorf("failed to find empty seats: %w", err)
	}
	if len(availableSeats) < seatCount {
		return nil, fmt.Errorf("seats were grabbed by another transaction, please try again")
	}

	// 3. Try to acquire Distributed Locks in Cache for these specific seats
	var lockedSeatIDs []string
	lockTTL := 3 * time.Minute

	for _, seat := range availableSeats {
		lockKey := fmt.Sprintf("seat_lock:%s", seat.ID)

		// Attempt to lock the seat in Redis
		acquired, err := s.cacheRepo.AcquireLock(ctx, lockKey, userID, lockTTL)
		if err != nil || !acquired {
			// ROLLBACK: If we fail to lock even one seat, we must release the ones we already locked!
			s.releaseLocks(ctx, lockedSeatIDs)
			return nil, fmt.Errorf("failed to acquire lock for seat %s, thundering herd collision", seat.ID)
		}
		lockedSeatIDs = append(lockedSeatIDs, seat.ID)
	}

	// 4. We successfully locked the seats! Now create a PENDING booking in the database
	bookingID := uuid.New().String()
	booking := &domain.Booking{
		ID:        bookingID,
		UserID:    userID,
		TrainID:   trainID,
		SeatIDs:   lockedSeatIDs,
		Status:    domain.BookingStatusPending,
		Amount:    float64(seatCount * 1500), // Assuming a flat rate for the flash sale
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.bookingRepo.Create(ctx, booking)
	if err != nil {
		// CRITICAL: If DB write fails, release the cache locks immediately
		s.releaseLocks(ctx, lockedSeatIDs)
		return nil, fmt.Errorf("failed to create pending booking: %w", err)
	}

	return booking, nil
}

// releaseLocks is a helper to clean up if a transaction fails halfway through
func (s *TatkalService) releaseLocks(ctx context.Context, seatIDs []string) {
	for _, id := range seatIDs {
		lockKey := fmt.Sprintf("seat_lock:%s", id)
		_ = s.cacheRepo.ReleaseLock(ctx, lockKey) // Fire and forget for cleanup
	}
}
