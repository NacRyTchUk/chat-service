package broker

import (
	"chat-service/internal/dto"
	"chat-service/pkg"
	"github.com/nats-io/nats.go"
	"log"
)

type MessageBroker interface {
	Publish(message dto.Message)
}

type messageBroker struct {
	broker *nats.Conn
}

func NewMessageBroker(broker *nats.Conn) MessageBroker {
	return &messageBroker{
		broker: broker,
	}
}

func (broker messageBroker) Publish(message dto.Message) {
	go func() {
		log.Println("broker publish call: ", message)
		bmsg, err := pkg.Serialize[dto.Message](message)
		if err != nil {
			return
		}
		_ = broker.broker.Publish("new-messages", bmsg)
		log.Println("broker publish complete: ", message)
	}()
}
