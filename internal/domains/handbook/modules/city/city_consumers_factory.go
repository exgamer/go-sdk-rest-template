package city

func NewCityConsumersFactory() *CityConsumersFactory {
	return &CityConsumersFactory{
		CityConsumer: NewCityConsumer(),
	}
}

type CityConsumersFactory struct {
	CityConsumer *CityConsumer
}
