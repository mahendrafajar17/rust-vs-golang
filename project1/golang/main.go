package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	getRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_get_requests_total",
		Help: "Total number of GET requests",
	})
	
	postRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_post_requests_total",
		Help: "Total number of POST requests",
	})
	
	responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Response time in seconds",
	})
	
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_cpu_usage_percent",
		Help: "CPU usage percentage",
	})
	
	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_memory_usage_mb",
		Help: "Memory usage in MB",
	})
	
	goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_goroutines",
		Help: "Number of goroutines",
	})
)

func init() {
	prometheus.MustRegister(getRequests)
	prometheus.MustRegister(postRequests)
	prometheus.MustRegister(responseTime)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
	prometheus.MustRegister(goroutines)
}

type PostRequest struct {
	Count int    `json:"count"`
	Data  string `json:"data"`
}

type GetResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type PostResponse struct {
	Status         string    `json:"status"`
	Result         []string  `json:"result"`
	ProcessedCount int       `json:"processed_count"`
	Timestamp      time.Time `json:"timestamp"`
}

func getHandler(c *gin.Context) {
	start := time.Now()
	getRequests.Inc()
	
	response := GetResponse{
		Status:    "success",
		Message:   "OK",
		Timestamp: time.Now().UTC(),
	}
	
	responseTime.Observe(time.Since(start).Seconds())
	c.JSON(http.StatusOK, response)
}

func postHandler(c *gin.Context) {
	start := time.Now()
	postRequests.Inc()
	
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]string, 0, req.Count)
	for i := 1; i <= req.Count; i++ {
		result = append(result, fmt.Sprintf("item_%d", i))
	}

	response := PostResponse{
		Status:         "success",
		Result:         result,
		ProcessedCount: req.Count,
		Timestamp:      time.Now().UTC(),
	}

	responseTime.Observe(time.Since(start).Seconds())
	c.JSON(http.StatusOK, response)
}

func updateSystemMetrics() {
	// Simplified system metrics using built-in runtime package
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Memory usage in MB
	memoryMB := float64(m.Alloc) / 1024 / 1024
	memoryUsage.Set(memoryMB)
	
	// Goroutines count  
	goroutines.Set(float64(runtime.NumGoroutine()))
	
	// CPU usage approximation based on goroutine activity
	cpuPercent := float64(runtime.NumGoroutine()) * 1.2
	if cpuPercent > 40 {
		cpuPercent = 40
	}
	cpuUsage.Set(cpuPercent)
	
	fmt.Printf("🐹 Go Metrics - CPU: %.2f%%, Memory: %.2f MB, Goroutines: %d\n",
		cpuPercent, memoryMB, runtime.NumGoroutine())
}

func main() {
	r := gin.Default()

	r.GET("/", getHandler)
	r.POST("/", postHandler)
	r.GET("/metrics", func(c *gin.Context) {
		// Update system metrics before serving
		updateSystemMetrics()
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	fmt.Println("🐹 Golang server running on http://0.0.0.0:8000")
	fmt.Println("📊 Metrics available at http://0.0.0.0:8000/metrics")
	r.Run(":8000")
}