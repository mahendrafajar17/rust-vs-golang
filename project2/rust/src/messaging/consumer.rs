use anyhow::Result;
use async_trait::async_trait;
use futures_util::stream::StreamExt;
use lapin::{
    message::Delivery,
    options::*,
    types::FieldTable,
    Channel,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Semaphore;
use tracing::{error, info, instrument, warn};
use uuid::Uuid;

use super::publisher::AMQPPublisher;
use crate::metrics::Metrics;

#[async_trait]
pub trait MessageHandler: Send + Sync {
    async fn handle(&self, delivery: Delivery) -> Result<Delivery>;
}

#[derive(Clone)]
pub struct AMQPConsumer {
    channel: Arc<Channel>,
    publisher: AMQPPublisher,
    metrics: Arc<Metrics>,
    concurrency: usize,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct InputMessage {
    pub user_id: String,
    pub product_name: String,
    pub quantity: i32,
    pub price: f64,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct OutputMessage {
    pub id: String,
    pub user_id: String,
    pub product_name: String,
    pub quantity: i32,
    pub price: f64,
}

#[derive(Clone)]
pub struct QueueProcessor {
    pub consumer: AMQPConsumer,
    pub input_queue: String,
    pub output_queue: String,
}

impl AMQPConsumer {
    pub fn new(channel: Channel, publisher: AMQPPublisher, metrics: Arc<Metrics>, concurrency: usize) -> Self {
        Self {
            channel: Arc::new(channel),
            publisher,
            metrics,
            concurrency,
        }
    }

    #[instrument(skip(self, handler))]
    pub async fn start_consuming<H>(&self, queue_name: &str, handler: H) -> Result<()>
    where
        H: MessageHandler + Clone + 'static,
    {
        // Update active consumers metric
        self.metrics.set_active_consumers(self.concurrency as f64);

        // Declare queue
        self.channel
            .queue_declare(
                queue_name,
                QueueDeclareOptions {
                    durable: true,
                    ..Default::default()
                },
                FieldTable::default(),
            )
            .await?;

        // Create consumer
        let mut consumer = self
            .channel
            .basic_consume(
                queue_name,
                "rust-consumer",
                BasicConsumeOptions::default(),
                FieldTable::default(),
            )
            .await?;

        info!(
            queue = queue_name,
            concurrency = self.concurrency,
            "Started consuming messages"
        );

        // Create semaphore to limit concurrency
        let semaphore = Arc::new(Semaphore::new(self.concurrency));

        // Process messages
        while let Some(delivery) = consumer.next().await {
            match delivery {
                Ok(delivery) => {
                    let handler = handler.clone();
                    let semaphore = semaphore.clone();
                    let metrics = self.metrics.clone();
                    let request_id = Uuid::new_v4().to_string();

                    tokio::spawn(async move {
                        let _permit = semaphore.acquire().await.unwrap();

                        let span = tracing::info_span!("message_processing", request_id = %request_id);
                        let _guard = span.enter();

                        // Increment received messages
                        metrics.inc_messages_received();

                        let start = std::time::Instant::now();

                        let result = handler.handle(delivery).await;
                        let duration = start.elapsed();
                        metrics.observe_processing_duration(duration);
                        
                        match result {
                            Ok(delivery) => {
                                metrics.inc_messages_processed();
                                
                                if let Err(e) = delivery.ack(BasicAckOptions::default()).await {
                                    error!(error = %e, "Failed to acknowledge message");
                                } else {
                                    info!(
                                        duration_ms = duration.as_millis(),
                                        "Message processed successfully"
                                    );
                                }
                            }
                            Err(e) => {
                                metrics.inc_messages_failed();
                                error!(error = %e, "Message processing failed");
                            }
                        }
                    });
                }
                Err(e) => {
                    error!(error = %e, "Failed to consume message");
                }
            }
        }

        warn!("Consumer stream ended");
        Ok(())
    }
}

impl QueueProcessor {
    pub fn new(consumer: AMQPConsumer, input_queue: String, output_queue: String) -> Self {
        Self {
            consumer,
            input_queue,
            output_queue,
        }
    }
}

#[async_trait]
impl MessageHandler for QueueProcessor {
    #[instrument(skip(self, delivery))]
    async fn handle(&self, delivery: Delivery) -> Result<Delivery> {
        // Parse input message
        let input_msg: InputMessage = serde_json::from_slice(&delivery.data)?;

        // Create output message with UUID
        let output_msg = OutputMessage {
            id: Uuid::new_v4().to_string(),
            user_id: input_msg.user_id.clone(),
            product_name: input_msg.product_name.clone(),
            quantity: input_msg.quantity,
            price: input_msg.price,
        };

        // Publish to output queue
        self.consumer
            .publisher
            .publish(&self.output_queue, &output_msg)
            .await?;

        info!(
            input_queue = %self.input_queue,
            output_queue = %self.output_queue,
            uuid_added = %output_msg.id,
            user_id = %output_msg.user_id,
            product_name = %output_msg.product_name,
            "Message forwarded with UUID"
        );

        Ok(delivery)
    }
}