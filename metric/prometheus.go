package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var MHttpRequestTotal *prometheus.CounterVec

type Prometheus struct {
	// Port is the network port where prometheus will serve an endpoint
	Port int
	// set true when system is under test
	test bool
}

// TODO: this doesn't work, at least not in basic testing.  We need to get a working Prometheus server, write some
// TODO: tests and add more metrics throughout the application.
// Run starts the Prometheus metrics server. Set metrics with the PrometheusMetrics.
func (p Prometheus) Run() {
	p.registerAllMetrics()

	path := "/metrics"
	log.Infof("starting prometheus at path '%s' on port '%d'", path, p.Port)
	if !p.test {
		mux := http.NewServeMux()
		mux.Handle(path, promhttp.Handler())
		go http.ListenAndServe(fmt.Sprintf(":%d", p.Port), promhttp.Handler())
	}
}

// registerAllMetrics ensures all custom metrics are individually defined and registered with Prometheus.
func (p Prometheus) registerAllMetrics() {
	MHttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of http requests to the server.",
		},
		[]string{"statuscode"},
	)
	prometheus.MustRegister(MHttpRequestTotal)
}
