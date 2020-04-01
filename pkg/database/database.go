package database

import (
	"database/sql"

	"github.com/jaceklubzinski/pd-checker/pkg/base"
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

//InitIncidentRepository create database schema
func (d *Store) InitIncidentRepository() {
	incidentsTable := `
	CREATE TABLE IF NOT EXISTS incidents(
		id TEXT NOT NULL,
		title TEXT NOT NULL,
		serviceid TEXT NOT NULL UNIQUE,
		createat TEXT NOT NULL,
		timer TEXT NOT NULL,
		alert TEXT DEFAULT "N",
		tocheck TEXT DEFAULT "N",
		trigger TEXT DEFAULT "N"
	);
	`
	servicesTable := `
	CREATE TABLE IF NOT EXISTS services(
		id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL
	);
	`
	d.createTable(incidentsTable)
	d.createTable(servicesTable)
}

func (d *Store) createTable(sqlTable string) {
	_, err := d.db.Exec(sqlTable)
	base.CheckErr(err)
}
