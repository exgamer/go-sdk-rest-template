package city

import (
	"github.com/exgamer/gosdk-rabbit-core/pkg/config"
)

func GetConsumers(
	consumersFactory *CityConsumersFactory,
) []config.HandlerRegister {
	return []config.HandlerRegister{
		{
			Handler: consumersFactory.CityConsumer.Consume,
			Config: config.NewConsumerTopicDurableConfig(
				"test-consumer",
				"test-rk",
				"test",
				"test",
				100,
			),
		},
	}
}
