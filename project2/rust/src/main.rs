mod amqp;
mod config;
mod messaging;
mod metrics;

use anyhow::Result;
use axum::{routing::get, Router};
use std::sync::Arc;
use tokio::{net::TcpListener, signal};
use tracing::{error, info, Level};
use tracing_subscriber::{fmt, prelude::*, EnvFilter};

use amqp::AMQPConnection;
use config::AppConfig;
use messaging::{AMQPConsumer, AMQPPublisher, QueueProcessor};
use metrics::Metrics;

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    let subscriber = tracing_subscriber::registry()
        .with(fmt::layer().json())
        .with(
            EnvFilter::builder()
                .with_default_directive(Level::INFO.into())
                .from_env_lossy(),
        );

    tracing::subscriber::set_global_default(subscriber)?;

    // Load configuration
    let config = AppConfig::load().map_err(|e| {
        error!("Failed to load configuration: {}", e);
        e
    })?;

    info!(
        app_name = %config.app.name,
        version = "1.0.0",
        "Starting application"
    );

    // Initialize metrics
    let app_metrics = Arc::new(Metrics::new());

    // Connect to RabbitMQ
    let connection = AMQPConnection::new(&config.amqp.url).await?;
    
    // Update connection metrics
    app_metrics.set_amqp_connections(1.0);

    // Create channels
    let publisher_channel = connection.create_channel().await?;
    let consumer_channel = connection
        .create_channel_with_qos(config.amqp.prefetch_count)
        .await?;

    // Setup publisher
    let publisher = AMQPPublisher::new(publisher_channel);

    // Setup consumer
    let consumer = AMQPConsumer::new(
        consumer_channel,
        publisher.clone(),
        app_metrics.clone(),
        config.amqp.concurrent,
    );

    // Setup queue processor
    let processor = QueueProcessor::new(
        consumer,
        config.queues.input_queue.clone(),
        config.queues.output_queue.clone(),
    );

    info!(
        input_queue = %config.queues.input_queue,
        output_queue = %config.queues.output_queue,
        concurrency = config.amqp.concurrent,
        http_port = config.app.port,
        "Queue processor started successfully"
    );

    // Start simple metrics server
    let metrics_clone = app_metrics.clone();
    let metrics_port = config.app.port;
    let metrics_handle = tokio::spawn(async move {
        let app = Router::new()
            .route("/metrics", get({
                let metrics = metrics_clone;
                || async move {
                    match metrics.encode_metrics() {
                        Ok(encoded) => encoded,
                        Err(e) => {
                            error!("Failed to encode metrics: {}", e);
                            "".to_string()
                        }
                    }
                }
            }));

        let addr = format!("0.0.0.0:{}", metrics_port);
        let listener = match TcpListener::bind(&addr).await {
            Ok(listener) => listener,
            Err(e) => {
                error!("Failed to bind metrics server: {}", e);
                return;
            }
        };
        
        info!("Metrics server listening on {}", addr);
        
        if let Err(e) = axum::serve(listener, app).await {
            error!("Metrics server error: {}", e);
        }
    });

    // Start consuming messages
    let consumer_handle = {
        let processor_clone = processor.clone();
        let consumer_clone = processor.consumer.clone();
        let input_queue = config.queues.input_queue.clone();
        tokio::spawn(async move {
            if let Err(e) = consumer_clone.start_consuming(&input_queue, processor_clone).await {
                error!("Consumer error: {}", e);
            }
        })
    };

    // Wait for shutdown signal
    info!("Application running. Press Ctrl+C to shutdown.");
    signal::ctrl_c().await?;
    info!("Shutdown signal received");

    // Graceful shutdown
    consumer_handle.abort();
    metrics_handle.abort();
    info!("Application shutdown complete");

    Ok(())
}