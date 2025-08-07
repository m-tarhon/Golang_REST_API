package main

import (
	"net/http"
	"time"
	"rest_api/metrics"
	"log"
)

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusCodeResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &statusCodeResponseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		metrics.ReqCounter.WithLabelValues(r.Method, http.StatusText(rw.statusCode), r.URL.Path).Inc()
		//metrics.ReqDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		log.Printf("Request: %s %s took %v", r.Method, r.URL.Path, duration)
	})
}