package city

import (
	"context"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

type Repository interface {
	Paginated(ctx context.Context, searchDto *Search) (*pagination.Paginated[City], error)
	GetById(ctx context.Context, id uint) (*City, error)
	Create(ctx context.Context, model *City) (*City, error)
	Update(ctx context.Context, model *City) error
	Delete(ctx context.Context, id uint) error
	Activate(ctx context.Context, id uint) error
	Deactivate(ctx context.Context, id uint) error
}

type HttpRepository interface {
	GetCity(ctx context.Context) (*City, error)
}

// CacheRepository - кеш поверх Repository. Реализация (Redis) живёт в
// Infrastructure и внедряется опционально: если её нет (например
// RedisKernel не зарегистрирован), Service работает напрямую с БД.
type CacheRepository interface {
	GetCityById(ctx context.Context, id uint) (*City, error)
	SetCity(ctx context.Context, model *City) error
}
