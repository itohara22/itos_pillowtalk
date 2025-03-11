package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sushi/handlers"
	"sushi/utils"
)

func main() {
	portCli := flag.String("port", ":6969", "port to serve on")
	flag.Parse() /// important dont forget

	db := utils.ConnectToDb()
	defer db.Close()

	tmp, err := parseTemplates()
	if err != nil {
		log.Panic(err.Error())
	}

	h := handlers.HandlerStruct{
		DbPool: db,
		T:      tmp,
	}

	ptToh := &h
	routerMux := router(ptToh)
	// we are using pointer to get struct at the memory address not the copy

	server := &http.Server{
		Addr:    *portCli,
		Handler: routerMux,
	}
	fmt.Printf("serving on %v\n", *portCli)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func parseTemplates() (*template.Template, error) {
	tmp, err := template.ParseGlob("templates/*.html")
	return tmp, err
}
