package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task_test/pkg/api"
	"task_test/pkg/db"
)

func main() {
	log.Print("Server has started")

	// Start the db
	pgdb, err := db.StartDB()
	if err != nil {
		log.Printf("error starting the database %v", err)
	}

	// Get the router of the API by passing the db
	router := api.StartAPI(pgdb)

	// Get the port from the environment variable
	port := os.Getenv("PORT")

	// Pass the router and start listening with the server
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
	}
}
