use axum::{
    extract::{Json, State},
    http::StatusCode,
    response::{Json as ResponseJson},
    routing::{get, post},
    Router,
};
use chrono::{DateTime, Utc};
use prometheus::{Counter, Encoder, Gauge, TextEncoder};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use sysinfo::{System, Pid};
use tokio::net::TcpListener;

#[derive(Clone)]
struct AppState {
    get_requests: Counter,
    post_requests: Counter,
    cpu_usage: Gauge,
    memory_usage: Gauge,
    system: Arc<std::sync::Mutex<System>>,
}

#[derive(Deserialize)]
struct PostRequest {
    count: usize,
    data: String,
}

#[derive(Serialize)]
struct GetResponse {
    status: String,
    message: String,
    timestamp: DateTime<Utc>,
}

#[derive(Serialize)]
struct PostResponse {
    status: String,
    result: Vec<String>,
    processed_count: usize,
    timestamp: DateTime<Utc>,
}

async fn get_handler(State(state): State<AppState>) -> Result<ResponseJson<GetResponse>, StatusCode> {
    state.get_requests.inc();
    
    let response = GetResponse {
        status: "success".to_string(),
        message: "OK".to_string(),
        timestamp: Utc::now(),
    };
    
    Ok(ResponseJson(response))
}

async fn post_handler(
    State(state): State<AppState>,
    Json(payload): Json<PostRequest>,
) -> Result<ResponseJson<PostResponse>, StatusCode> {
    state.post_requests.inc();
    
    let mut result = Vec::new();
    
    for i in 1..=payload.count {
        result.push(format!("item_{}", i));
    }
    
    let response = PostResponse {
        status: "success".to_string(),
        result,
        processed_count: payload.count,
        timestamp: Utc::now(),
    };
    
    Ok(ResponseJson(response))
}

async fn metrics_handler(State(state): State<AppState>) -> Result<String, StatusCode> {
    println!("Metrics handler called!");
    
    // Update system metrics
    if let Ok(mut sys) = state.system.lock() {
        sys.refresh_all();
        
        // Get current process info
        let current_pid = Pid::from(std::process::id() as usize);
        if let Some(process) = sys.process(current_pid) {
            let cpu_usage = process.cpu_usage() as f64;
            let memory_usage = process.memory() as f64 / 1024.0 / 1024.0; // Convert to MB
            
            state.cpu_usage.set(cpu_usage);
            state.memory_usage.set(memory_usage);
            
            println!("🦀 Rust Metrics - CPU: {:.2}%, Memory: {:.2} MB", cpu_usage, memory_usage);
        }
    }
    
    let encoder = TextEncoder::new();
    let metric_families = prometheus::gather();
    
    println!("Metrics families count: {}", metric_families.len());
    
    let mut buffer = Vec::new();
    match encoder.encode(&metric_families, &mut buffer) {
        Ok(()) => {
            let metrics_string = String::from_utf8_lossy(&buffer).to_string();
            println!("Metrics output length: {}", metrics_string.len());
            Ok(metrics_string)
        },
        Err(e) => {
            println!("Error encoding metrics: {:?}", e);
            Err(StatusCode::INTERNAL_SERVER_ERROR)
        }
    }
}

#[tokio::main]
async fn main() {
    println!("Starting Rust server with simple metrics...");
    
    // Create metrics using global registry
    let get_requests = Counter::new("http_get_requests_total", "Total number of GET requests").unwrap();
    let post_requests = Counter::new("http_post_requests_total", "Total number of POST requests").unwrap();
    let cpu_usage = Gauge::new("process_cpu_usage_percent", "CPU usage percentage").unwrap();
    let memory_usage = Gauge::new("process_memory_usage_mb", "Memory usage in MB").unwrap();
    
    // Register metrics with global registry
    prometheus::register(Box::new(get_requests.clone())).unwrap();
    prometheus::register(Box::new(post_requests.clone())).unwrap();
    prometheus::register(Box::new(cpu_usage.clone())).unwrap();
    prometheus::register(Box::new(memory_usage.clone())).unwrap();
    
    let system = Arc::new(std::sync::Mutex::new(System::new_all()));
    
    let state = AppState {
        get_requests,
        post_requests,
        cpu_usage,
        memory_usage,
        system,
    };
    
    let app = Router::new()
        .route("/", get(get_handler))
        .route("/", post(post_handler))
        .route("/metrics", get(metrics_handler))
        .with_state(state);

    let listener = TcpListener::bind("127.0.0.1:3002").await.unwrap();
    println!("Rust server running on http://127.0.0.1:3002");
    println!("Metrics available at http://127.0.0.1:3002/metrics");
    
    axum::serve(listener, app).await.unwrap();
}