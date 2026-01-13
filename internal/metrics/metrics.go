package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Total number of probe checks
var ProbeTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "probe_requests_total",
		Help: "Total number of probe requests",
	},
	[]string{"target"},
)

// Total number of probe checks that failed
var ProbeFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "probe_failed",
	Help: "Total number of probe checks that failed",
},
	[]string{"target"},
)

// Latency of probe checks (in seconds)
var ProbeLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "probe_latency_seconds",
	Help: "Latency of probe checks in seconds",
},
	[]string{"target"},
)

func Register() {
	prometheus.MustRegister(ProbeTotal, ProbeFailures, ProbeLatency)
}
