use prometheus::{
    Encoder, Gauge, GaugeVec, Histogram, HistogramOpts, IntCounter, Opts, Registry,
    TextEncoder,
};
use sysinfo::{Pid, System};
use std::sync::{Arc, Mutex};
use tracing::debug;

#[derive(Clone)]
pub struct Metrics {
    // Message processing metrics
    pub messages_received: IntCounter,
    pub messages_processed: IntCounter,
    pub messages_failed: IntCounter,
    pub processing_duration: Histogram,
    
    // Queue metrics
    pub queue_depth: GaugeVec,
    pub active_consumers: Gauge,
    
    // System metrics
    pub cpu_usage: Gauge,
    pub memory_usage: Gauge,
    
    // Connection metrics
    pub amqp_connections: Gauge,
    pub amqp_reconnections: IntCounter,
    
    // System info
    system: Arc<Mutex<System>>,
    registry: Registry,
}

impl Metrics {
    pub fn new() -> Self {
        let registry = Registry::new();

        let messages_received = IntCounter::with_opts(Opts::new(
            "rabbitmq_messages_received_total",
            "Total number of messages received from RabbitMQ",
        )).unwrap();

        let messages_processed = IntCounter::with_opts(Opts::new(
            "rabbitmq_messages_processed_total", 
            "Total number of messages successfully processed",
        )).unwrap();

        let messages_failed = IntCounter::with_opts(Opts::new(
            "rabbitmq_messages_failed_total",
            "Total number of messages that failed processing",
        )).unwrap();

        let processing_duration = Histogram::with_opts(HistogramOpts::new(
            "rabbitmq_message_processing_seconds",
            "Time taken to process a message",
        ).buckets(vec![0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5]))
        .unwrap();

        let queue_depth = GaugeVec::new(
            Opts::new("rabbitmq_queue_depth", "Number of messages in queue"),
            &["queue_name"],
        ).unwrap();

        let active_consumers = Gauge::with_opts(Opts::new(
            "rabbitmq_active_consumers",
            "Number of active consumer workers",
        )).unwrap();

        let cpu_usage = Gauge::with_opts(Opts::new(
            "process_cpu_usage_percent",
            "CPU usage percentage",
        )).unwrap();

        let memory_usage = Gauge::with_opts(Opts::new(
            "process_memory_usage_mb",
            "Memory usage in MB",
        )).unwrap();

        let amqp_connections = Gauge::with_opts(Opts::new(
            "amqp_connections_active",
            "Number of active AMQP connections",
        )).unwrap();

        let amqp_reconnections = IntCounter::with_opts(Opts::new(
            "amqp_reconnections_total",
            "Total number of AMQP reconnections",
        )).unwrap();

        // Register all metrics
        registry.register(Box::new(messages_received.clone())).unwrap();
        registry.register(Box::new(messages_processed.clone())).unwrap();
        registry.register(Box::new(messages_failed.clone())).unwrap();
        registry.register(Box::new(processing_duration.clone())).unwrap();
        registry.register(Box::new(queue_depth.clone())).unwrap();
        registry.register(Box::new(active_consumers.clone())).unwrap();
        registry.register(Box::new(cpu_usage.clone())).unwrap();
        registry.register(Box::new(memory_usage.clone())).unwrap();
        registry.register(Box::new(amqp_connections.clone())).unwrap();
        registry.register(Box::new(amqp_reconnections.clone())).unwrap();

        Self {
            messages_received,
            messages_processed,
            messages_failed,
            processing_duration,
            queue_depth,
            active_consumers,
            cpu_usage,
            memory_usage,
            amqp_connections,
            amqp_reconnections,
            system: Arc::new(Mutex::new(System::new_all())),
            registry,
        }
    }

    pub fn update_system_metrics(&self) {
        if let Ok(mut system) = self.system.lock() {
            system.refresh_all();
            
            let current_pid = Pid::from(std::process::id() as usize);
            if let Some(process) = system.process(current_pid) {
                let cpu_usage = process.cpu_usage() as f64;
                let memory_usage = process.memory() as f64 / 1024.0 / 1024.0; // Convert to MB
                
                self.cpu_usage.set(cpu_usage);
                self.memory_usage.set(memory_usage);
                
                debug!(
                    cpu_percent = cpu_usage,
                    memory_mb = memory_usage,
                    "System metrics updated"
                );
            }
        }
    }

    pub fn gather_metrics(&self) -> Vec<prometheus::proto::MetricFamily> {
        self.update_system_metrics();
        self.registry.gather()
    }

    pub fn encode_metrics(&self) -> Result<String, Box<dyn std::error::Error>> {
        let encoder = TextEncoder::new();
        let metric_families = self.gather_metrics();
        
        let mut buffer = Vec::new();
        encoder.encode(&metric_families, &mut buffer)?;
        
        Ok(String::from_utf8(buffer)?)
    }

    // Helper methods for easy metric updates
    pub fn inc_messages_received(&self) {
        self.messages_received.inc();
    }

    pub fn inc_messages_processed(&self) {
        self.messages_processed.inc();
    }

    pub fn inc_messages_failed(&self) {
        self.messages_failed.inc();
    }

    pub fn observe_processing_duration(&self, duration: std::time::Duration) {
        self.processing_duration.observe(duration.as_secs_f64());
    }

    pub fn set_queue_depth(&self, queue_name: &str, depth: f64) {
        self.queue_depth.with_label_values(&[queue_name]).set(depth);
    }

    pub fn set_active_consumers(&self, count: f64) {
        self.active_consumers.set(count);
    }

    pub fn set_amqp_connections(&self, count: f64) {
        self.amqp_connections.set(count);
    }

    pub fn inc_amqp_reconnections(&self) {
        self.amqp_reconnections.inc();
    }
}