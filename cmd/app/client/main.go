package main

import (
	"fmt"
	"os"

	"github.com/Liar233/rabbitmq-stub/internal/app/client"
)

func main() {

	if len(os.Args) < 2 {

		_, _ = fmt.Fprintln(os.Stderr, "Error: wrong cli argument")
		os.Exit(1)

		return
	}

	config := &client.AppClientConfig{
		Dsl: "amqp://user:secret@rabbitmq:5672/",
		Config: &client.PublisherConfig{
			ExchangeName:       "test_exchange_name",
			ExchangeType:       "direct",
			ExchangeDurable:    true,
			ExchangeAutoDelete: false,
		},
	}

	appClient := client.NewAppClient(config)

	if err := appClient.Bootstrap(); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, "Bootstrapping error:", err.Error())
		os.Exit(1)
	}

	if err := appClient.Fire(os.Args[1], "test_exchange_test_queue"); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, "Fire error:", err.Error())
		os.Exit(1)
	}
}
