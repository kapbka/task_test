package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"

	"task_test/pkg/db/models"
)

type MetricsResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Metrics []*models.Metric `json:"metrics"`
}

// start api with the pgdb and return a chi router
func StartAPI(pgdb *pg.DB) *chi.Mux {
	// Get the router
	r := chi.NewRouter()

	// Add middleware
	// In this case we will store our DB to use it later
	r.Use(middleware.Logger)

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		getMetrics(pgdb, w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}

func getMetrics(pgdb *pg.DB, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Retrieve fromInt and to_ts query parameter values
	fromInt, errFrom := strconv.ParseInt(queryParams.Get("from_ts"), 10, 64)
	toInt, errTo := strconv.ParseInt(queryParams.Get("to_ts"), 10, 64)
	if errFrom != nil || errTo != nil {
		res := &MetricsResponse{
			Success: false,
			Error:   "Timestamp range is not valid!",
			Metrics: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call models package to access the database and return the metrics
	fromTs := time.Unix(fromInt*1000, 0)
	toTs := time.Unix(toInt*1000, 0)
	metrics, err := models.GetMetrics(pgdb, fromTs, toTs)
	if err != nil {
		res := &MetricsResponse{
			Success: false,
			Error:   err.Error(),
			Metrics: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Positive response
	res := &MetricsResponse{
		Success: true,
		Error:   "",
		Metrics: metrics,
	}

	// Encode the positive response to json and send it back
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding metrics: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
