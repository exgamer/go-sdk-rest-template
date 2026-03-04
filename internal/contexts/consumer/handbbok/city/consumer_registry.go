package city

import (
	"github.com/exgamer/gosdk-rabbit-core/pkg/config"
)

func GetConsumers(
	consumer *Consumer,
) []config.HandlerRegister {
	return []config.HandlerRegister{
		{
			Handler: consumer.Consume,
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
