package database

import (
	"database/sql"

	"github.com/jaceklubzinski/pd-checker/pkg/config"
)

type Store struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) *Store {
	return &Store{db}
}
func ConnectDatabase(config *config.Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}
