package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"project2-golang/config"
	"project2-golang/provider/amqpx"
	"project2-golang/provider/messaging"
	"project2-golang/provider/metrics"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".", "./config")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.WithFields(logrus.Fields{
		"app_name": cfg.App.Name,
		"version":  "1.0.0",
	}).Info("Starting application")

	// Initialize metrics
	appMetrics := metrics.NewMetrics()
	
	// Connect to RabbitMQ
	conn, err := amqpx.Dial(cfg.AMQP.GetDSN())
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to RabbitMQ")
	}
	defer conn.Close()
	
	// Update connection metrics
	appMetrics.SetAMQPConnections(1)

	// Create channel for publisher
	publisherChannel, err := conn.Channel()
	if err != nil {
		logger.WithError(err).Fatal("Failed to create publisher channel")
	}
	defer publisherChannel.Close()

	// Create channel for consumer
	consumerChannel, err := conn.Channel()
	if err != nil {
		logger.WithError(err).Fatal("Failed to create consumer channel")
	}
	defer consumerChannel.Close()

	// Setup publisher
	publisher := messaging.NewAMQPPublisher(publisherChannel)

	// Setup consumer
	consumer := messaging.NewAMQPConsumer(
		consumerChannel,
		publisher,
		appMetrics,
		cfg.AMQP.Concurrent,
		cfg.AMQP.PrefetchCount,
	)

	// Setup queue processor
	processor := messaging.NewQueueProcessor(
		consumer,
		cfg.Queues.InputQueue,
		cfg.Queues.OutputQueue,
	)

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start metrics collection
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				appMetrics.UpdateSystemMetrics()
			}
		}
	}()

	// Start simple metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.WithField("port", cfg.App.Port).Info("Starting metrics server")
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.App.Port), nil); err != nil {
			logger.WithError(err).Error("Metrics server error")
		}
	}()

	// Start consuming messages
	err = consumer.StartConsuming(ctx, cfg.Queues.InputQueue, processor)
	if err != nil {
		logger.WithError(err).Fatal("Failed to start consumer")
	}

	logger.WithFields(logrus.Fields{
		"input_queue":  cfg.Queues.InputQueue,
		"output_queue": cfg.Queues.OutputQueue,
		"concurrency":  cfg.AMQP.Concurrent,
	}).Info("Queue processor started successfully")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := consumer.Stop(shutdownCtx); err != nil {
		logger.WithError(err).Error("Error during consumer shutdown")
	}

	cancel() // Cancel main context
	logger.Info("Application shutdown complete")
}