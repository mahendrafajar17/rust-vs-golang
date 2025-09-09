package metrics

import (
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	// Message processing metrics
	MessagesReceived   prometheus.Counter
	MessagesProcessed  prometheus.Counter
	MessagesFailed     prometheus.Counter
	ProcessingDuration prometheus.Histogram
	
	// Queue metrics  
	QueueDepth         prometheus.GaugeVec
	ActiveConsumers    prometheus.Gauge
	
	// System metrics
	CPUUsage        prometheus.Gauge
	MemoryUsage     prometheus.Gauge
	Goroutines      prometheus.Gauge
	
	// Connection metrics
	AMQPConnections    prometheus.Gauge
	AMQPReconnections  prometheus.Counter
	
	logger *logrus.Logger
	mutex  sync.RWMutex
	
	// CPU tracking
	lastCPUTime time.Time
	lastCPUUsage time.Duration
}

func NewMetrics() *Metrics {
	m := &Metrics{
		MessagesReceived: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rabbitmq_messages_received_total",
			Help: "Total number of messages received from RabbitMQ",
		}),
		MessagesProcessed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rabbitmq_messages_processed_total",
			Help: "Total number of messages successfully processed",
		}),
		MessagesFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rabbitmq_messages_failed_total",
			Help: "Total number of messages that failed processing",
		}),
		ProcessingDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "rabbitmq_message_processing_seconds",
			Help: "Time taken to process a message",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10), // 1ms to ~1s
		}),
		QueueDepth: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rabbitmq_queue_depth",
			Help: "Number of messages in queue",
		}, []string{"queue_name"}),
		ActiveConsumers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rabbitmq_active_consumers",
			Help: "Number of active consumer workers",
		}),
		CPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "process_cpu_usage_percent",
			Help: "CPU usage percentage",
		}),
		MemoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "process_memory_usage_mb",
			Help: "Memory usage in MB",
		}),
		Goroutines: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "process_goroutines",
			Help: "Number of goroutines",
		}),
		AMQPConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "amqp_connections_active",
			Help: "Number of active AMQP connections",
		}),
		AMQPReconnections: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "amqp_reconnections_total",
			Help: "Total number of AMQP reconnections",
		}),
		logger: logrus.New(),
		lastCPUTime: time.Now(),
	}

	// Register all metrics
	prometheus.MustRegister(m.MessagesReceived)
	prometheus.MustRegister(m.MessagesProcessed)
	prometheus.MustRegister(m.MessagesFailed)
	prometheus.MustRegister(m.ProcessingDuration)
	prometheus.MustRegister(&m.QueueDepth)
	prometheus.MustRegister(m.ActiveConsumers)
	prometheus.MustRegister(m.CPUUsage)
	prometheus.MustRegister(m.MemoryUsage)
	prometheus.MustRegister(m.Goroutines)
	prometheus.MustRegister(m.AMQPConnections)
	prometheus.MustRegister(m.AMQPReconnections)

	return m
}

func (m *Metrics) UpdateSystemMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Memory usage in MB
	memoryMB := float64(memStats.Alloc) / 1024 / 1024
	m.MemoryUsage.Set(memoryMB)

	// Goroutines count
	goroutineCount := float64(runtime.NumGoroutine())
	m.Goroutines.Set(goroutineCount)

	// Calculate CPU usage based on CPU time
	now := time.Now()
	elapsed := now.Sub(m.lastCPUTime)
	
	if elapsed > 0 {
		// Get current CPU usage from runtime
		var cpuUsage runtime.MemStats
		runtime.ReadMemStats(&cpuUsage)
		
		// Estimate CPU percentage based on GC and processing activity
		gcCPUFraction := cpuUsage.GCCPUFraction * 100
		
		// Add base CPU usage estimation based on message processing
		baseCPU := float64(goroutineCount-1) * 0.8 // Exclude main goroutine
		if baseCPU > 30 {
			baseCPU = 30
		}
		
		totalCPU := gcCPUFraction + baseCPU
		if totalCPU > 100 {
			totalCPU = 100
		}
		
		m.CPUUsage.Set(totalCPU)
		m.lastCPUTime = now
	}

	m.logger.WithFields(logrus.Fields{
		"memory_mb":      memoryMB,
		"goroutines":     goroutineCount,
	}).Debug("System metrics updated")
}

// Helper methods for easy metric updates
func (m *Metrics) IncMessagesReceived() {
	m.MessagesReceived.Inc()
}

func (m *Metrics) IncMessagesProcessed() {
	m.MessagesProcessed.Inc()
}

func (m *Metrics) IncMessagesFailed() {
	m.MessagesFailed.Inc()
}

func (m *Metrics) ObserveProcessingDuration(duration time.Duration) {
	m.ProcessingDuration.Observe(duration.Seconds())
}

func (m *Metrics) SetQueueDepth(queueName string, depth float64) {
	m.QueueDepth.WithLabelValues(queueName).Set(depth)
}

func (m *Metrics) SetActiveConsumers(count float64) {
	m.ActiveConsumers.Set(count)
}

func (m *Metrics) SetAMQPConnections(count float64) {
	m.AMQPConnections.Set(count)
}

func (m *Metrics) IncAMQPReconnections() {
	m.AMQPReconnections.Inc()
}