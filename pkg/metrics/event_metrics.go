package metrics

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type RecordMetricsEvent struct {
	Event *prometheus.CounterVec
}

func (e *RecordMetricsEvent) NewRecordMetricsEvent() {
	e.Event = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pd_event_connection_total",
			Help: "The total number of connection to PagerDuty API events",
		},
		[]string{"code"},
	)
	prometheus.MustRegister(e.Event)
}

func (e *RecordMetricsEvent) RecordMetricsEventError(toParse string) {
	r, _ := regexp.Compile("HTTP Status Code: ([0-9]+), Message")
	match := r.FindStringSubmatch(toParse)[1]
	match40x, _ := regexp.Compile("^40[0-9]$")
	match50x, _ := regexp.Compile("^50[0-9]$")
	switch {
	case match40x.MatchString(match):
		e.Event.With(prometheus.Labels{"code": "40x"}).Inc()

	case match50x.MatchString(match):
		e.Event.With(prometheus.Labels{"code": "50x"}).Inc()
	}
}

func (e *RecordMetricsEvent) RecordMetricsEvent(toParse string) {
	match, _ := regexp.MatchString("success", toParse)
	if match {
		e.Event.With(prometheus.Labels{"code": "20x"}).Inc()
	}
}
