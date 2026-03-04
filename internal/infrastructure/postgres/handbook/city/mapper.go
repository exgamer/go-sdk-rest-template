package city

import citydomain "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"

// modelToEntity преобразует DB модель в доменную сущность
func modelToEntity(model *city) *citydomain.City {
	if model == nil {
		return nil
	}

	return &citydomain.City{
		ID:     model.ID,
		Name:   model.Name,
		Status: model.Status,
	}
}

// entityToModel преобразует доменную сущность в DB модель
func entityToModel(entity *citydomain.City) *city {
	if entity == nil {
		return nil
	}

	return &city{
		ID:     entity.ID,
		Name:   entity.Name,
		Status: entity.Status,
	}
}
