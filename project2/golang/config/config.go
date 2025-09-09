package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	AMQP   AMQPConfig   `mapstructure:"amqp"`
	Queues QueueConfig  `mapstructure:"queues"`
	Log    LogConfig    `mapstructure:"logging"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type AMQPConfig struct {
	Scheme       string `mapstructure:"scheme"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Concurrent   int    `mapstructure:"concurrent"`
	PrefetchCount int   `mapstructure:"prefetch_count"`
	PrefetchSize  int   `mapstructure:"prefetch_size"`
	Global       bool   `mapstructure:"global"`
}

type QueueConfig struct {
	InputQueue  string `mapstructure:"input_queue"`
	OutputQueue string `mapstructure:"output_queue"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig(paths ...string) (*Config, error) {
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *AMQPConfig) GetDSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", 
		c.Scheme, c.Username, c.Password, c.Host, c.Port)
}