package client

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherConfig struct {
	ExchangeName       string
	ExchangeType       string
	ExchangeDurable    bool
	ExchangeAutoDelete bool
}

type Publisher struct {
	client *Client
	config *PublisherConfig
}

func (publisher *Publisher) Fire(ctx context.Context, routingKey string, payload []byte) error {

	if err := publisher.DeclareExchange(); err != nil {

		return err
	}

	msg := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    amqp.Persistent,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Now(),
		Type:            "",
		UserId:          "",
		AppId:           "rabbitmq-publisher",
		Body:            payload,
	}

	return publisher.client.channel.PublishWithContext(
		ctx,
		publisher.config.ExchangeName,
		routingKey,
		false,
		false,
		msg,
	)
}

func (publisher *Publisher) DeclareExchange() error {

	err := publisher.client.channel.ExchangeDeclare(
		publisher.config.ExchangeName,
		publisher.config.ExchangeType,
		publisher.config.ExchangeDurable,
		publisher.config.ExchangeAutoDelete,
		false,
		false,
		nil,
	)

	if err != nil {

		return err
	}

	return nil
}

func NewPublisher(client *Client, config *PublisherConfig) *Publisher {

	return &Publisher{
		client: client,
		config: config,
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
