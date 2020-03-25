package incident

import (
	"fmt"

	"github.com/jaceklubzinski/pd-checker/pkg/database"
)

type Server struct {
	Incident *IncidentService
	Database *database.Store
}

func NewServer(Incident *IncidentService, Database *database.Store) *Server {
	return &Server{Incident: Incident, Database: Database}
}

func (s *Server) GetIncident() {
	testinc := s.Database.GetIncident()
	for _, v := range testinc {
		fmt.Println(v)
	}
}
