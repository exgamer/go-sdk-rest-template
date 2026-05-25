package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/entrypoint/base/rabbit/handbbok/city"
)

func newConsumersFactory() *consumersFactory {
	return &consumersFactory{
		CityConsumer: city.NewConsumer(),
	}
}

type consumersFactory struct {
	CityConsumer *city.Consumer
}
