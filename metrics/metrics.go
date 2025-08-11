package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ill need a prometheus server running
// and then a grafaa server to visualize the metrics, running in front o fit
var (
	Registry       = prometheus.NewRegistry()
	SystemRegistry = prometheus.NewRegistry()

	ReqCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "status", "endpoint"})

	CpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "Current CPU usage",
	}, []string{"cpu"})

	MemoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "Current Memory usage",
	}, []string{"memory"})

	// ReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	// 	Name:    "http_request_duration_seconds",
	// 	Help:    "Duration of HTTP requests in seconds",
	// 	Buckets: prometheus.DefBuckets,
	// }, []string{"method", "endpoint"})

	FileOpCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "file_operations_total",
		Help: "Total number of file operations",
	}, []string{"operation", "file_type", "status"})

	FileOpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "file_operation_duration_seconds",
		Help:    "Duration of file operations in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation", "file_type"})
)

func Init() {
	Registry.MustRegister(ReqCounter)
	Registry.MustRegister(FileOpCounter)
	Registry.MustRegister(FileOpDuration)
	// Registry.MustRegister(ReqDuration)
	Registry.MustRegister(CpuUsage) // to add cpu usage i need an external package
	SystemRegistry.MustRegister(MemoryUsage)
	go collectSystemMetrics()
}

var SystemMetricsInterval = 10 * time.Second

func collectSystemMetrics() {
	ticker := time.NewTicker(SystemMetricsInterval)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		cpu := runtime.NumCPU()

		MemoryUsage.WithLabelValues("alloc").Set(float64(m.Alloc))
		MemoryUsage.WithLabelValues("sys").Set(float64(m.Sys))
		MemoryUsage.WithLabelValues("heap_alloc").Set(float64(m.HeapAlloc))
		MemoryUsage.WithLabelValues("heap_sys").Set(float64(m.HeapSys))
		CpuUsage.WithLabelValues("usage").Set(float64(cpu))
	}
}

// Wrapper function for timing file operations
func MeasureFileOperation(operation, fileType string, fn func() error) error {
	start := time.Now()

	err := fn()

	duration := time.Since(start).Seconds()
	FileOpDuration.WithLabelValues(operation, fileType).Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}
	FileOpCounter.WithLabelValues(operation, fileType, status).Inc()

	return err
}
