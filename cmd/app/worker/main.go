package main

import (
	"fmt"
	"os"

	"github.com/Liar233/rabbitmq-stub/internal/app/worker"
)

func main() {

	config := &worker.AppWorkerConfig{
		Dsl: "amqp://user:secret@rabbitmq:5672/",
		Config: &worker.ConsumerConfig{
			ConsumerName:       "test_consumer_name",
			AutoAck:            false,
			ExchangeName:       "test_exchange_name",
			ExchangeType:       "direct",
			ExchangeDurable:    true,
			ExchangeAutoDelete: false,
			QueueName:          "test_queue_name",
			QueueDurable:       true,
			QueueAutoDelete:    false,
			BindingKey:         "test_exchange_test_queue",
		},
	}

	appWorker := worker.NewAppWorker(config)

	if err := appWorker.Bootstrap(); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, "Bootstrapping error:", err.Error())
		os.Exit(1)
	}

	if err := appWorker.Run(); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, "Running error:", err.Error())
		os.Exit(1)
	}
}
