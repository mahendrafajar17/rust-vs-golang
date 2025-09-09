use serde::{Deserialize, Serialize};
use config::{Config, ConfigError, File};

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct AppConfig {
    pub app: App,
    pub amqp: Amqp,
    pub queues: Queues,
    pub logging: Logging,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct App {
    pub name: String,
    pub port: u16,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Amqp {
    pub url: String,
    pub concurrent: usize,
    pub prefetch_count: u16,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Queues {
    pub input_queue: String,
    pub output_queue: String,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Logging {
    pub level: String,
    pub format: String,
}

impl AppConfig {
    pub fn load() -> Result<Self, ConfigError> {
        let settings = Config::builder()
            .add_source(File::with_name("config"))
            .add_source(File::with_name("config.yaml").required(false))
            .build()?;

        settings.try_deserialize()
    }
}