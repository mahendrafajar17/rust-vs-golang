use anyhow::Result;
use lapin::{
    options::*,
    types::FieldTable,
    BasicProperties, Channel,
};
use serde::Serialize;
use std::sync::Arc;
use tracing::{info, instrument};

#[derive(Clone)]
pub struct AMQPPublisher {
    channel: Arc<Channel>,
}

impl AMQPPublisher {
    pub fn new(channel: Channel) -> Self {
        Self {
            channel: Arc::new(channel),
        }
    }

    #[instrument(skip(self, message))]
    pub async fn publish<T>(&self, queue_name: &str, message: &T) -> Result<()>
    where
        T: Serialize,
    {
        // Serialize message to JSON
        let payload = serde_json::to_vec(message)?;

        // Declare queue to ensure it exists
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

        // Publish message
        self.channel
            .basic_publish(
                "",
                queue_name,
                BasicPublishOptions::default(),
                &payload,
                BasicProperties::default()
                    .with_content_type("application/json".into())
                    .with_delivery_mode(2), // Persistent message
            )
            .await?
            .await?; // Wait for confirmation

        info!(
            queue = queue_name,
            message_size = payload.len(),
            "Message published successfully"
        );

        Ok(())
    }

    #[instrument(skip(self, message))]
    pub async fn publish_with_routing_key<T>(
        &self,
        exchange: &str,
        routing_key: &str,
        message: &T,
    ) -> Result<()>
    where
        T: Serialize,
    {
        let payload = serde_json::to_vec(message)?;

        self.channel
            .basic_publish(
                exchange,
                routing_key,
                BasicPublishOptions::default(),
                &payload,
                BasicProperties::default()
                    .with_content_type("application/json".into())
                    .with_delivery_mode(2),
            )
            .await?
            .await?;

        info!(
            exchange = exchange,
            routing_key = routing_key,
            message_size = payload.len(),
            "Message published to exchange"
        );

        Ok(())
    }
}