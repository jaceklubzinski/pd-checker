package database

import (
	"database/sql"

	"github.com/jaceklubzinski/pd-checker/pkg/config"
)

type store struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) *store {
	return &store{db}
}
func ConnectDatabase(config *config.Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}
