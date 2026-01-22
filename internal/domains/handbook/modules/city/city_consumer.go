package city

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewCityConsumer() *CityConsumer {
	return &CityConsumer{}
}

type CityConsumer struct {
}

func (c *CityConsumer) Consume(ctx context.Context, msg *message.Message) error {
	fmt.Println(string(msg.Payload))

	return nil
}
