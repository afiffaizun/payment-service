package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"payment-service/internal/delivery"
	"payment-service/internal/repository"
	"payment-service/internal/usecase"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user_payment")
	dbPassword := getEnv("DB_PASSWORD", "pass_payment")
	dbName := getEnv("DB_NAME", "db_payment")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	repo := repository.NewPostgresRepo(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := delivery.NewHttpHandler(uc)

	r := chi.NewRouter()

	// Tambahkan health check endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment Service is running"))
	})

	r.Post("/transfer", handler.Transfer)
	r.Post("/topup", handler.TopUp)
	r.Get("/transaction/{refId}", handler.GetTransaction)
	r.Get("/wallet/{userId}", handler.GetWallet)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
