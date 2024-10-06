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
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	r.Route("/metrics", func(r chi.Router) {
		r.Get("/", getMetrics)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Retrieve from_int and to_ts query parameter values
	from_int, err_from := strconv.ParseInt(queryParams.Get("from_ts"), 10, 64)
	to_int, err_to := strconv.ParseInt(queryParams.Get("to_ts"), 10, 64)
	if err_from != nil || err_to != nil {
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

	// Get db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &MetricsResponse{
			Success: false,
			Error:   "could not get DB from context",
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
	from_ts := time.Unix(from_int*1000, 0)
	to_ts := time.Unix(to_int*1000, 0)
	metrics, err := models.GetMetrics(pgdb, from_ts, to_ts)
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
