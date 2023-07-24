package action

import (
	"fmt"
	"math/rand"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageProcessor struct {
}

func (m *MessageProcessor) Process(msg amqp.Delivery) error {

	fmt.Println("New message:")
	fmt.Println("")
	fmt.Println("ContentType :", msg.ContentType)
	fmt.Println("ContentEncoding :", msg.ContentEncoding)
	fmt.Println("DeliveryMode :", msg.DeliveryMode)
	fmt.Println("Priority :", msg.Priority)
	fmt.Println("CorrelationId :", msg.CorrelationId)
	fmt.Println("ReplyTo :", msg.ReplyTo)
	fmt.Println("Expiration :", msg.Expiration)
	fmt.Println("MessageId :", msg.MessageId)
	fmt.Println("Timestamp :", msg.Timestamp)
	fmt.Println("Type :", msg.Type)
	fmt.Println("UserId :", msg.UserId)
	fmt.Println("AppId :", msg.AppId)
	fmt.Println("ConsumerTag :", msg.ConsumerTag)
	fmt.Println("DeliveryTag :", msg.DeliveryTag)
	fmt.Println("Redelivered :", msg.Redelivered)
	fmt.Println("Exchange :", msg.Exchange)
	fmt.Println("RoutingKey :", msg.RoutingKey)
	fmt.Println("Body :", string(msg.Body))
	fmt.Println()

	r := rand.Int()

	if r%3 == 0 {

		return fmt.Errorf("artificial processing error")
	}

	return nil
}

func NewMessageProcessor() *MessageProcessor {

	return &MessageProcessor{}
}
