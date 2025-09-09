package messaging

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"project2-golang/provider/amqpx"
	"project2-golang/provider/metrics"
)

type AMQPConsumer struct {
	channel     *amqpx.Channel
	publisher   *AMQPPublisher
	metrics     *metrics.Metrics
	logger      *logrus.Logger
	wg          *sync.WaitGroup
	state       int32
	concurrency int
}

type MessageHandler interface {
	Handle(ctx context.Context, data amqp.Delivery) error
}

type QueueProcessor struct {
	consumer    *AMQPConsumer
	inputQueue  string
	outputQueue string
	logger      *logrus.Logger
}

// Input message structure
type InputMessage struct {
	UserID      string  `json:"user_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// Output message structure (with added UUID)
type OutputMessage struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

func NewAMQPConsumer(channel *amqpx.Channel, publisher *AMQPPublisher, metrics *metrics.Metrics, concurrency int) *AMQPConsumer {
	return &AMQPConsumer{
		channel:     channel,
		publisher:   publisher,
		metrics:     metrics,
		logger:      logrus.New(),
		wg:          &sync.WaitGroup{},
		concurrency: concurrency,
	}
}

func NewQueueProcessor(consumer *AMQPConsumer, inputQueue, outputQueue string) *QueueProcessor {
	return &QueueProcessor{
		consumer:    consumer,
		inputQueue:  inputQueue,
		outputQueue: outputQueue,
		logger:      logrus.New(),
	}
}

func (c *AMQPConsumer) StartConsuming(ctx context.Context, queue string, handler MessageHandler) error {
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return nil // Already running
	}

	// Update active consumers metric
	c.metrics.SetActiveConsumers(float64(c.concurrency))

	// Set QoS
	err := c.channel.Qos(50, 0, false)
	if err != nil {
		return err
	}

	// Declare queue
	_, err = c.channel.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Start consumer goroutines
	for i := 0; i < c.concurrency; i++ {
		c.wg.Add(1)
		go c.worker(ctx, queue, handler, i)
	}

	c.logger.WithFields(logrus.Fields{
		"queue":       queue,
		"concurrency": c.concurrency,
	}).Info("Started AMQP consumer")

	return nil
}

func (c *AMQPConsumer) worker(ctx context.Context, queue string, handler MessageHandler, workerID int) {
	defer c.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			c.logger.WithFields(logrus.Fields{
				"worker_id": workerID,
				"panic":     r,
			}).Error("Consumer worker panicked")
		}
	}()

	messages, err := c.channel.Consume(
		queue,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		c.logger.WithError(err).WithField("worker_id", workerID).Error("Failed to register consumer")
		return
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.WithField("worker_id", workerID).Info("Consumer worker stopping")
			return
		case delivery, ok := <-messages:
			if !ok {
				c.logger.WithField("worker_id", workerID).Info("Consumer channel closed")
				return
			}

			c.processMessage(ctx, delivery, handler, workerID)
		}
	}
}

func (c *AMQPConsumer) processMessage(ctx context.Context, delivery amqp.Delivery, handler MessageHandler, workerID int) {
	requestID := uuid.New().String()
	msgCtx := context.WithValue(ctx, "request_id", requestID)

	logger := c.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"worker_id":  workerID,
		"queue":      delivery.RoutingKey,
	})

	// Increment received messages
	c.metrics.IncMessagesReceived()
	
	logger.Info("Processing message")

	start := time.Now()
	err := handler.Handle(msgCtx, delivery)
	duration := time.Since(start)

	// Record processing duration
	c.metrics.ObserveProcessingDuration(duration)

	if err != nil {
		c.metrics.IncMessagesFailed()
		logger.WithError(err).WithField("duration_ms", duration.Milliseconds()).Error("Message processing failed")
		delivery.Nack(false, true) // Requeue message
	} else {
		c.metrics.IncMessagesProcessed()
		logger.WithField("duration_ms", duration.Milliseconds()).Info("Message processed successfully")
		delivery.Ack(false)
	}
}

func (c *AMQPConsumer) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&c.state, 1, 0) {
		return nil // Already stopped
	}

	c.logger.Info("Stopping AMQP consumer")
	
	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.logger.Info("All consumer workers stopped")
	case <-time.After(30 * time.Second):
		c.logger.Warn("Consumer stop timeout reached")
	}

	return nil
}

// QueueProcessor implementation
func (qp *QueueProcessor) Handle(ctx context.Context, delivery amqp.Delivery) error {
	// Parse input message
	var inputMsg InputMessage
	if err := json.Unmarshal(delivery.Body, &inputMsg); err != nil {
		return err
	}

	// Add UUID and create output message
	outputMsg := OutputMessage{
		ID:          uuid.New().String(),
		UserID:      inputMsg.UserID,
		ProductName: inputMsg.ProductName,
		Quantity:    inputMsg.Quantity,
		Price:       inputMsg.Price,
	}

	// Publish to output queue
	err := qp.consumer.publisher.Publish(ctx, qp.outputQueue, outputMsg)
	if err != nil {
		return err
	}

	qp.logger.WithFields(logrus.Fields{
		"request_id":    ctx.Value("request_id"),
		"input_queue":   qp.inputQueue,
		"output_queue":  qp.outputQueue,
		"uuid_added":    outputMsg.ID,
		"user_id":       outputMsg.UserID,
		"product_name":  outputMsg.ProductName,
	}).Info("Message forwarded with UUID")

	return nil
}