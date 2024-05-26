package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"sync"
)

// Pool struct that holds RabbitMQ connections and channels
type Pool struct {
	mu             sync.Mutex
	connections    []*amqp.Connection
	channels       []*amqp.Channel
	connURL        string
	maxConn        int
	maxChanPerConn int
}

// NewPool initializes a new Pool
func NewPool(maxConn, maxChanPerConn int) (*Pool, error) {

	connURL := os.Getenv("RABBITMQ_URL")

	pool := &Pool{
		connURL:        connURL,
		maxConn:        maxConn,
		maxChanPerConn: maxChanPerConn,
	}
	if err := pool.init(); err != nil {
		return nil, err
	}
	return pool, nil
}

// init initializes the connection pool
func (p *Pool) init() error {
	for i := 0; i < p.maxConn; i++ {
		conn, err := amqp.Dial(p.connURL)
		if err != nil {
			return err
		}
		p.connections = append(p.connections, conn)
		for j := 0; j < p.maxChanPerConn; j++ {
			ch, err := conn.Channel()
			if err != nil {
				return err
			}

			_, err = ch.QueueDeclare(
				"seeds",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("Failed to declare a queue: %s", err.Error())
			}

			p.channels = append(p.channels, ch)

		}
	}
	return nil
}

// GetChannel retrieves a channel from the pool
func (p *Pool) GetChannel() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.channels) == 0 {
		return nil, amqp.ErrClosed
	}

	ch := p.channels[0]
	p.channels = p.channels[1:]
	return ch, nil
}

// ReturnChannel returns a channel to the pool
func (p *Pool) ReturnChannel(ch *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.channels = append(p.channels, ch)
}

// Close closes all connections and channels in the pool
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ch := range p.channels {
		if err := ch.Close(); err != nil {
			return err
		}
	}

	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	p.channels = nil
	p.connections = nil
	return nil
}

// Publish publishes a message to the given exchange and routing key
func (p *Pool) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ch, err := p.GetChannel()
	if err != nil {
		return err
	}
	defer p.ReturnChannel(ch)

	err = ch.Publish(
		exchange,  // exchange
		key,       // routing key
		mandatory, // mandatory
		immediate, // immediate
		msg,       // message to publish
	)
	return err
}

// Consume sets up a consumer for the queue and processes messages
func (p *Pool) Consume(consumerName string, handler func(amqp.Delivery)) error {
	ch, err := p.GetChannel()
	if err != nil {
		return err
	}
	defer p.ReturnChannel(ch)

	msgs, err := ch.Consume(
		"seeds",      // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	// Start a goroutine to handle messages
	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	return nil
}
