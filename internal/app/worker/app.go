package worker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Liar233/rabbitmq-stub/internal/app/worker/action"
)

type AppWorkerConfig struct {
	Dsl    string
	Config *ConsumerConfig
}

type AppWorker struct {
	config    *AppWorkerConfig
	consumer  *Consumer
	client    *Client
	processor AmqpProcessorInterface
}

func (app *AppWorker) Bootstrap() error {

	var err error

	if app.client, err = NewClient(app.config.Dsl); err != nil {

		return err
	}

	return nil
}

func (app *AppWorker) Run() error {

	app.processor = action.NewMessageProcessor()

	app.consumer = NewConsumer(app.client, app.config.Config, app.processor)

	if err := app.consumer.DeclareExchange(); err != nil {

		return err
	}

	if err := app.consumer.DeclareQueue(); err != nil {

		return err
	}

	if err := app.consumer.DeclareBinding(); err != nil {

		return err
	}

	go app.Stop()

	return app.consumer.Run()
}

func (app *AppWorker) Stop() {

	sigintChan := make(chan os.Signal, 1)

	defer close(sigintChan)

	signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigintChan

	_, _ = fmt.Fprintf(os.Stderr, "\ngot signal: %s\n", s)

	app.consumer.ShutDown()
}

func NewAppWorker(config *AppWorkerConfig) *AppWorker {

	return &AppWorker{config: config}
}
