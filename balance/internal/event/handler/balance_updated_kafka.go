package handler

import (
	"balances/pkg/events"
	"balances/pkg/kafka"
	"fmt"
	"sync"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Consumer
}

func NewBalanceUpdatedKafkaHandler(k *kafka.Consumer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		Kafka: k,
	}
}

func (b *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	// b.Kafka.Publish(message, nil, "balances")
	b.Kafka.Consume()
	fmt.Println("BalanceUpdatedKafkaHandler: ", message.GetPayload())
}
