use anyhow::Result;
use lapin::{
    options::*,
    types::FieldTable,
    Connection, ConnectionProperties,
};
use std::sync::Arc;
use tracing::{error, info, warn};

#[derive(Clone)]
pub struct AMQPConnection {
    connection: Arc<Connection>,
}

impl AMQPConnection {
    pub async fn new(amqp_url: &str) -> Result<Self> {
        let connection = Self::connect_with_retry(amqp_url).await?;
        
        Ok(Self {
            connection: Arc::new(connection),
        })
    }

    async fn connect_with_retry(amqp_url: &str) -> Result<Connection> {
        let mut attempts = 0;
        let max_attempts = 5;

        loop {
            match Connection::connect(amqp_url, ConnectionProperties::default()).await {
                Ok(connection) => {
                    info!("Successfully connected to RabbitMQ");
                    return Ok(connection);
                }
                Err(e) => {
                    attempts += 1;
                    if attempts >= max_attempts {
                        error!("Failed to connect to RabbitMQ after {} attempts: {}", max_attempts, e);
                        return Err(e.into());
                    }
                    
                    warn!("Failed to connect to RabbitMQ (attempt {}): {}. Retrying...", attempts, e);
                    tokio::time::sleep(tokio::time::Duration::from_secs(3)).await;
                }
            }
        }
    }

    pub async fn create_channel(&self) -> Result<lapin::Channel> {
        let channel = self.connection.create_channel().await?;
        Ok(channel)
    }

    pub async fn create_channel_with_qos(&self, prefetch_count: u16) -> Result<lapin::Channel> {
        let channel = self.connection.create_channel().await?;
        
        channel
            .basic_qos(prefetch_count, BasicQosOptions::default())
            .await?;
            
        Ok(channel)
    }

    pub fn is_connected(&self) -> bool {
        self.connection.status().connected()
    }
}