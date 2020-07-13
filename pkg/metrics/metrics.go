package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

//Server prometheus metrics server
func Server(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info(http.ListenAndServe(":"+port, nil))
}
