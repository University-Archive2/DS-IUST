package repository

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type AggregatedMetricsRepo interface {
	Create(a *AggregatedMetric) error
	BatchCreate(a []*AggregatedMetric) error
}

type aggregatedMetricsRepo struct {
	db *sqlx.DB
}

type AggregatedMetric struct {
	Symbol    string    `db:"symbol"`
	Metric    string    `db:"metric"`
	Value     float64   `db:"value"`
	Timestamp time.Time `db:"timestamp"`
	CreatedAt string    `db:"created_at"`
}

var schema = `
CREATE TABLE IF NOT EXISTS aggregated_metrics (
	symbol VARCHAR(255) NOT NULL,
	metric VARCHAR(255) NOT NULL,
	value FLOAT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (symbol, metric, timestamp)
);
`

func NewAggregatedMetricsRepo(db *sqlx.DB) AggregatedMetricsRepo {
	db.MustExec(schema)
	return &aggregatedMetricsRepo{
		db: db,
	}
}

func (r *aggregatedMetricsRepo) Create(a *AggregatedMetric) error {
	_, err := r.db.NamedExec("INSERT INTO aggregated_metrics (metric, value, timestamp) VALUES (:metric, :value, :timestamp)", a)
	return err
}

func (r *aggregatedMetricsRepo) BatchCreate(a []*AggregatedMetric) error {
	_, err := r.db.NamedExec("INSERT INTO aggregated_metrics (symbol, metric, value, timestamp) VALUES (:symbol, :metric, :value, :timestamp)", a)
	return err
}
