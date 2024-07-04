package contracts

type QueueService interface {
	Publish(BaseMessage) error
	Consume(handlersMap map[string]MessageHandler) error
	Close()
}

type BaseMessage interface {
	SetType()
	GetType() string
}

type MessageHandler interface {
	HandleMessage(messageBody map[string]interface{}) error
}
