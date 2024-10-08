package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Metric struct {
	Ts          time.Time `json:"ts"`
	CpuLoad     float64   `json:"cpu_load"`
	Concurrency int64     `json:"concurrency"`
}

func InsertMetric(db *pg.DB, req *Metric) (bool, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return false, err
	}

	return true, err
}

func GetMetrics(db *pg.DB, fromTs time.Time, toTs time.Time) ([]*Metric, error) {
	// Query filtered results
	metrics := make([]*Metric, 0)
	err := db.Model(&metrics).
		Where("ts >= ?", fromTs).
		Where("ts <= ?", toTs).
		Select()

	return metrics, err
}
