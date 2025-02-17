package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sushi/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connString = "postgresql://dante:password@localhost:5432/dante"
)

func main() {
	portCli := flag.String("port", ":6969", "port to serve on")
	flag.Parse() /// important dont forget

	db := connectToDb()
	defer db.Close()

	h := handlers.HandlerStruct{
		DbPool: db,
	}
	ptToh := &h
	routerMux := router(ptToh)
	// we are using pointer to get struct at the memory address not the copy

	server := &http.Server{
		Addr:    *portCli,
		Handler: routerMux,
	}
	fmt.Printf("serving on %v\n", *portCli)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func connectToDb() *pgxpool.Pool {
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
