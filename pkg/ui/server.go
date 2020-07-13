package ui

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Server) Listen() {
	http.HandleFunc("/incidents", s.incident)
	log.Info(http.ListenAndServe(":9090", nil))
}

func (s Server) incident(w http.ResponseWriter, req *http.Request) {
	incidents, err := s.DbRepository.GetIncident()
	if err != nil {
		panic(err)
	}
	t := template.Must(template.ParseFiles("pkg/ui/incidents.tmpl"))
	err = t.Execute(w, incidents)
	if err != nil {
		panic(err)
	}
}
