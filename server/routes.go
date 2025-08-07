package main

import (
	"net/http"
	"rest_api/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes() *http.ServeMux {
		mux := http.NewServeMux()
		
		mux.HandleFunc("/healthcheck", MetricsMiddleware(http.HandlerFunc(healthcheck)).ServeHTTP)
		mux.Handle("/users", MetricsMiddleware(basicAuth(userManagement)))
		mux.Handle("/users/", MetricsMiddleware(basicAuth(userManagement)))
		mux.Handle("/apps", MetricsMiddleware(http.HandlerFunc(appsManagement)))
		mux.Handle("/apps/", MetricsMiddleware(http.HandlerFunc(appsManagement)))

		// uses default registry, but i have a custom one
		mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))

		return mux
}
