package utils

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	connString = "postgresql://dante:password@localhost:5432/dante"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedBytes), err
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ConnectToDb() *pgxpool.Pool {
	// conte, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	return dbPool
}

func RunMigrations(db *pgxpool.Pool) {
	migrationQuery := `
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			director TEXT NOT NULL,
			rating DOUBLE PRECISION CHECK (rating >= 0 AND rating <= 10)
		);
	`

	_, err := db.Exec(context.Background(), migrationQuery)
	if err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	log.Println("Migrations applied successfully")
}
