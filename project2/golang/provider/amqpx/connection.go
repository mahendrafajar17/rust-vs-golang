package amqpx

import (
	"sync"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Connection struct {
	conn     *amqp.Connection
	channels map[*Channel]struct{}
	mutex    sync.RWMutex
	url      string
	closed   int32
	logger   *logrus.Logger
}

type Channel struct {
	ch     *amqp.Channel
	conn   *Connection
	closed int32
	qos    *QoSSettings
	mutex  sync.RWMutex
}

type QoSSettings struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

func Dial(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	c := &Connection{
		conn:     conn,
		channels: make(map[*Channel]struct{}),
		url:      url,
		logger:   logrus.New(),
	}

	go c.handleReconnect()
	return c, nil
}

func (c *Connection) Channel() (*Channel, error) {
	if atomic.LoadInt32(&c.closed) == 1 {
		return nil, amqp.ErrClosed
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		ch:   ch,
		conn: c,
	}

	c.mutex.Lock()
	c.channels[channel] = struct{}{}
	c.mutex.Unlock()

	return channel, nil
}

func (c *Connection) handleReconnect() {
	notifyClose := c.conn.NotifyClose(make(chan *amqp.Error, 1))
	
	for {
		select {
		case err := <-notifyClose:
			if err != nil {
				c.logger.WithError(err).Error("AMQP connection lost, attempting to reconnect")
				c.reconnect()
			}
		}
	}
}

func (c *Connection) reconnect() {
	for {
		time.Sleep(3 * time.Second)
		
		conn, err := amqp.Dial(c.url)
		if err != nil {
			c.logger.WithError(err).Error("Failed to reconnect to AMQP")
			continue
		}

		c.mutex.Lock()
		c.conn = conn
		
		// Recreate all channels
		for channel := range c.channels {
			newCh, err := conn.Channel()
			if err != nil {
				c.logger.WithError(err).Error("Failed to recreate channel")
				continue
			}
			
			channel.mutex.Lock()
			channel.ch = newCh
			
			// Restore QoS settings if they exist
			if channel.qos != nil {
				err = newCh.Qos(
					channel.qos.PrefetchCount,
					channel.qos.PrefetchSize,
					channel.qos.Global,
				)
				if err != nil {
					c.logger.WithError(err).Error("Failed to restore QoS settings")
				}
			}
			
			channel.mutex.Unlock()
		}
		c.mutex.Unlock()

		c.logger.Info("Successfully reconnected to AMQP")
		
		// Start new reconnect handler (fixed goroutine leak)
		go c.handleReconnect()
		
		break
	}
}

func (c *Connection) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return amqp.ErrClosed
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for channel := range c.channels {
		channel.Close()
	}

	return c.conn.Close()
}

// Channel methods
func (ch *Channel) Qos(prefetchCount, prefetchSize int, global bool) error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	err := ch.ch.Qos(prefetchCount, prefetchSize, global)
	if err == nil {
		ch.qos = &QoSSettings{
			PrefetchCount: prefetchCount,
			PrefetchSize:  prefetchSize,
			Global:        global,
		}
	}
	return err
}

func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	return ch.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	return ch.ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	return ch.ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (ch *Channel) Close() error {
	if !atomic.CompareAndSwapInt32(&ch.closed, 0, 1) {
		return amqp.ErrClosed
	}

	ch.conn.mutex.Lock()
	delete(ch.conn.channels, ch)
	ch.conn.mutex.Unlock()

	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	
	return ch.ch.Close()
}