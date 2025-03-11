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
