package messaging

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"project2-golang/provider/amqpx"
)

type AMQPPublisher struct {
	pool    *sync.Pool
	channel *amqpx.Channel
}

type AMQPPublisherOptions struct {
	Exchange    string
	RoutingKey  string
	Mandatory   bool
	Immediate   bool
	ContentType string
	Priority    uint8
	Persistent  bool
}

func NewAMQPPublisher(channel *amqpx.Channel) *AMQPPublisher {
	return &AMQPPublisher{
		channel: channel,
		pool: &sync.Pool{
			New: func() interface{} {
				return &AMQPPublisherOptions{
					ContentType: "application/json",
					Persistent:  true,
				}
			},
		},
	}
}

func (p *AMQPPublisher) Publish(ctx context.Context, queue string, message interface{}, options ...func(*AMQPPublisherOptions)) error {
	opts := p.pool.Get().(*AMQPPublisherOptions)
	defer p.pool.Put(opts)

	// Reset options
	*opts = AMQPPublisherOptions{
		ContentType: "application/json",
		Persistent:  true,
	}

	// Apply custom options
	for _, option := range options {
		option(opts)
	}

	// Marshal message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Declare queue to ensure it exists
	_, err = p.channel.QueueDeclare(
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

	// Prepare AMQP message
	msg := amqp.Publishing{
		ContentType:  opts.ContentType,
		Body:         body,
		Timestamp:    time.Now(),
		Priority:     opts.Priority,
	}

	if opts.Persistent {
		msg.DeliveryMode = amqp.Persistent
	}

	// Publish message
	return p.channel.Publish(
		opts.Exchange,
		queue, // routing key (queue name for direct exchange)
		opts.Mandatory,
		opts.Immediate,
		msg,
	)
}

// Option functions
func WithExchange(exchange string) func(*AMQPPublisherOptions) {
	return func(opts *AMQPPublisherOptions) {
		opts.Exchange = exchange
	}
}

func WithRoutingKey(routingKey string) func(*AMQPPublisherOptions) {
	return func(opts *AMQPPublisherOptions) {
		opts.RoutingKey = routingKey
	}
}

func WithPriority(priority uint8) func(*AMQPPublisherOptions) {
	return func(opts *AMQPPublisherOptions) {
		opts.Priority = priority
	}
}

func WithMandatory(mandatory bool) func(*AMQPPublisherOptions) {
	return func(opts *AMQPPublisherOptions) {
		opts.Mandatory = mandatory
	}
}