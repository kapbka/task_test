package main

import (
	"log"
	"math/rand"
	"task_test/pkg/db"
	"task_test/pkg/db/models"
	"time"
)

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)

}

func main() {
	// Start DB
	pgdb, err := db.StartDB()
	if err != nil {
		log.Printf("error starting the database %v", err)
	}

	// Get the current time
	now := time.Now()

	// Calculate the time 5 minutes ago
	fiveMinutesAgo := now.Add(-5 * time.Minute)

	// Define the interval (60 milliseconds)
	interval := 60 * time.Millisecond

	// Generate timestamps from 5 minutes ago to now
	for t := fiveMinutesAgo; t.Before(now) || t.Equal(now); t = t.Add(interval) {
		_, err = models.InsertMetric(pgdb, &models.Metric{
			Ts:          t,
			CpuLoad:     randFloats(0, 100),
			Concurrency: int64(rand.Intn(500000)),
		})
		if err != nil {
			log.Printf("error starting the database %v", err)
		}
	}

}
