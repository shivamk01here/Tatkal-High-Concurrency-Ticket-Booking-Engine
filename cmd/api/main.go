package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/shivamk01here/tatkal-engine/internal/adapters/cache/redis"
	"github.com/shivamk01here/tatkal-engine/internal/adapters/db/postgres"
	handler "github.com/shivamk01here/tatkal-engine/internal/adapters/handlers/http"
	"github.com/shivamk01here/tatkal-engine/internal/core/services"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	trainRepo := postgres.NewTrainRepository(dbPool)
	seatRepo := postgres.NewSeatRepository(dbPool)
	cacheRepo := redisadapter.NewCacheRepository(redisClient)

	tatkalService := services.NewTatkalService(trainRepo, seatRepo, nil, cacheRepo)

	bookingHandler := handler.NewBookingHandler(tatkalService)

	r := chi.NewRouter()
	r.Post("/api/v1/bookings", bookingHandler.HandleInitiateBooking)

	log.Println("Tatkal Engine starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
