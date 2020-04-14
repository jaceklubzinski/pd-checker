package ui

import "github.com/jaceklubzinski/pd-checker/pkg/database"

type Server struct {
	DbRepository database.IncidentRepository
}

func NewServer(DbRepository database.IncidentRepository) *Server {
	return &Server{DbRepository: DbRepository}
}
