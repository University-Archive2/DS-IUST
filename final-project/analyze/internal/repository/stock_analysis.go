package repository

import "github.com/jmoiron/sqlx"

type StockAnalysisRepo interface {
	Create(a *StockAnalysis) error
}

type stockAnalysisRepo struct {
	db *sqlx.DB
}

type StockAnalysis struct {
	Symbol    string  `db:"symbol"`
	Indicator string  `db:"indicator"`
	Value     float64 `db:"value"`
	Timestamp float64 `db:"timestamp"`
	CreatedAt string  `db:"created_at"`
}

var schema = `
CREATE TABLE IF NOT EXISTS stock_analysis (
	symbol VARCHAR(255) NOT NULL,
	indicator VARCHAR(255) NOT NULL,
	value FLOAT NOT NULL,
	timestamp FLOAT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (symbol, indicator, value, timestamp)
);
`

func NewStockAnalysisRepo(db *sqlx.DB) StockAnalysisRepo {
	db.MustExec(schema)
	return &stockAnalysisRepo{
		db: db,
	}
}

func (r *stockAnalysisRepo) Create(a *StockAnalysis) error {
	_, err := r.db.NamedExec("INSERT INTO stock_analysis (symbol, indicator, value, timestamp) VALUES (:symbol, :indicator, :value, :timestamp)", a)
	return err
}
