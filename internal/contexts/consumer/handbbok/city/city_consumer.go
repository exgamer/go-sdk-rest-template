package city

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewConsumer() *Consumer {
	return &Consumer{}
}

type Consumer struct {
}

func (c *Consumer) Consume(ctx context.Context, msg *message.Message) error {
	fmt.Println(string(msg.Payload))

	return nil
}
