package client

import (
	"context"
)

type AppClientConfig struct {
	Dsl    string
	Config *PublisherConfig
}

type AppClient struct {
	Dsl    string
	config *AppClientConfig
	client *Client
}

func (app *AppClient) Bootstrap() error {

	var err error

	if app.client, err = NewClient(app.config.Dsl); err != nil {

		return err
	}

	return nil
}

func (app *AppClient) Fire(msg string, routingKey string) error {

	publisher := NewPublisher(app.client, app.config.Config)

	if err := publisher.DeclareExchange(); err != nil {

		return err
	}

	ctx := context.Background()

	return publisher.Fire(ctx, routingKey, []byte(msg))
}

func NewAppClient(config *AppClientConfig) *AppClient {

	return &AppClient{config: config}
}
