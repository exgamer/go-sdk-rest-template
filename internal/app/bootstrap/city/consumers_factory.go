package city

import "github.com/exgamer/go-sdk-rest-template/internal/transport/consumer/handbbok/city"

func newConsumersFactory() *consumersFactory {
	return &consumersFactory{
		CityConsumer: city.NewConsumer(),
	}
}

type consumersFactory struct {
	CityConsumer *city.Consumer
}
