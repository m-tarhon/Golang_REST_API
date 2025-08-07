package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var(
	Registry = prometheus.NewRegistry()
	ReqCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "status", "endpoint"})
)

func Init() {
	Registry.MustRegister(ReqCounter)
}

func SetupMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

}
