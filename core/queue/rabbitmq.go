package queue

import (
	"encoding/json"
	"gotranslate/core/contracts"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQ(url string, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: ch, queue: queueName}, nil
}

var _ contracts.QueueService = (*RabbitMQ)(nil)

func (r *RabbitMQ) Publish(data contracts.BaseMessage) error {
	data.SetType()
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	return r.channel.Publish("", r.queue, false, false, msg)
}

func (r *RabbitMQ) Consume(handlersMap map[string]contracts.MessageHandler) error {
	msgs, err := r.channel.Consume(r.queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var data map[string]interface{}
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("unable to unmarshal body to consume message on queue %v. Error: %v\n", r.queue, err.Error())
				return
			}

			typeValue, exists := data["Type"]
			if !exists {
				log.Println("error parsing message, no Type found")
			}

			messageType, ok := typeValue.(string)
			if !ok {
				log.Println("unable to read Type of message")
			}

			handler, found := handlersMap[messageType]
			if !found {
				log.Printf("message of type %v can't be handled, make sure it's registered", messageType)
			}
			handler.HandleMessage(data)
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
