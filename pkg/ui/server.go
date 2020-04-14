package ui

import (
	"html/template"
	"net/http"
)

func (s Server) Listen() {
	http.HandleFunc("/incidents", s.incident)
	http.ListenAndServe(":9090", nil)
}

func (s Server) incident(w http.ResponseWriter, req *http.Request) {
	incidents := s.DbRepository.GetIncident()
	t := template.Must(template.ParseFiles("pkg/ui/incidents.tmpl"))
	err := t.Execute(w, incidents)
	if err != nil {
		panic(err)
	}
}
