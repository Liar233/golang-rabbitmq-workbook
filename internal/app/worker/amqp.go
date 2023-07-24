package worker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpProcessorInterface interface {
	Process(msg amqp.Delivery) error
}

type ConsumerConfig struct {
	ConsumerName string
	AutoAck      bool

	ExchangeName       string
	ExchangeType       string
	ExchangeDurable    bool
	ExchangeAutoDelete bool

	QueueName       string
	QueueDurable    bool
	QueueAutoDelete bool

	BindingName string
	BindingKey  string
}

type Consumer struct {
	client    *Client
	config    *ConsumerConfig
	processor AmqpProcessorInterface
	queue     *amqp.Queue
	stop      chan struct{}
}

func (consumer *Consumer) Run() error {

	defer func() {
		_ = consumer.client.conn.Close()
		_ = consumer.client.channel.Close()
	}()

	delivery, err := consumer.client.channel.Consume(
		consumer.config.QueueName,
		consumer.config.ConsumerName,
		consumer.config.AutoAck,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		return err
	}

	deliverErr := make(chan error)

	go consumer.handler(delivery, deliverErr)

	return <-deliverErr
}

func (consumer *Consumer) handler(deliveries <-chan amqp.Delivery, errCh chan error) {

	for {
		select {
		case msg, ok := <-deliveries:

			if !ok {

				// Channel is closed
			}

			if err := consumer.processor.Process(msg); err != nil {

				errCh <- err
				return
			}

			if err := msg.Ack(true); err != nil {

				errCh <- err
				return
			}

		case <-consumer.stop:

			errCh <- nil
			return
		}
	}
}

func (consumer *Consumer) ShutDown() {

	consumer.stop <- struct{}{}
}

func (consumer *Consumer) DeclareExchange() error {

	err := consumer.client.channel.ExchangeDeclare(
		consumer.config.ExchangeName,
		consumer.config.ExchangeType,
		consumer.config.ExchangeDurable,
		consumer.config.ExchangeAutoDelete,
		false,
		false,
		nil,
	)

	if err != nil {

		return err
	}

	return nil
}

func (consumer *Consumer) DeclareQueue() error {

	queue, err := consumer.client.channel.QueueDeclare(
		consumer.config.QueueName,
		consumer.config.QueueDurable,
		consumer.config.QueueAutoDelete,
		false,
		false,
		nil,
	)

	if err != nil {

		return err
	}

	consumer.queue = &queue

	return err
}

func (consumer *Consumer) DeclareBinding() error {

	err := consumer.client.channel.QueueBind(
		consumer.config.QueueName,
		consumer.config.BindingKey,
		consumer.config.ExchangeName,
		false,
		nil,
	)

	if err != nil {

		return err
	}

	return err
}

func NewConsumer(
	client *Client,
	config *ConsumerConfig,
	processor AmqpProcessorInterface,
) *Consumer {

	return &Consumer{
		client:    client,
		config:    config,
		processor: processor,
		stop:      make(chan struct{}),
	}
}

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (c *Client) Close() error {

	if err := c.channel.Close(); err != nil {

		return err
	}

	if err := c.conn.Close(); err != nil {

		return err
	}

	return nil
}

func NewClient(dsl string) (*Client, error) {

	conn, err := amqp.Dial(dsl)

	if err != nil {

		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {

		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}
