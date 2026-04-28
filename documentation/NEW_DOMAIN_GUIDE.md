# Пошаговое создание нового модуля

Пример: модуль `product` в домене `catalog`.

Пути: `domains/catalog/product/`, `transport/.../catalog/product/`, `infrastructure/.../catalog/product/`, `bootstrap/product/`.

## Содержание

1. [Domain Layer](#1-domain-layer)
2. [Infrastructure — Postgres](#2-infrastructure--postgres)
3. [Infrastructure — Redis](#3-infrastructure--redis)
4. [Infrastructure — HTTP Client](#4-infrastructure--http-client)
5. [Transport — HTTP Handler](#5-transport--http-handler)
6. [Transport — RabbitMQ Consumer](#6-transport--rabbitmq-consumer)
7. [Bootstrap — Module](#7-bootstrap--module)
8. [Регистрация в App](#8-регистрация-в-app)

---

## 1. Domain Layer

**`internal/domains/catalog/product/entity.go`**

```go
package product

type Product struct {
    ID         uint
    Name       string
    Price      float64
    CategoryID uint
    Status     int
}
```

> Никаких GORM-тегов, никаких JSON-тегов. Только бизнес-поля.

---

**`internal/domains/catalog/product/search.go`**

```go
package product

type Search struct {
    ID         uint
    Name       string
    CategoryID uint
    Page       uint
    PerPage    uint
}
```

---

**`internal/domains/catalog/product/repository.go`**

```go
package product

import (
    "context"

    "github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

type Repository interface {
    Paginated(ctx context.Context, s *Search) (*pagination.Paginated[Product], error)
    GetById(ctx context.Context, id uint) (*Product, error)
    Create(ctx context.Context, m *Product) (*Product, error)
    Update(ctx context.Context, m *Product) error
    Delete(ctx context.Context, id uint) error
    Activate(ctx context.Context, id uint) error
    Deactivate(ctx context.Context, id uint) error
}
```

> Domain определяет контракт. Infrastructure реализует его.

---

**`internal/domains/catalog/product/service.go`**

```go
package product

import (
    "context"

    "github.com/exgamer/gosdk-core/pkg/debug"
    "github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

type Service struct {
    repository Repository
}

func NewService(repository Repository) *Service {
    return &Service{repository: repository}
}

func (s *Service) Paginated(ctx context.Context, search *Search) (*pagination.Paginated[Product], error) {
    return s.repository.Paginated(ctx, search)
}

func (s *Service) GetById(ctx context.Context, id uint) (*Product, error) {
    if dbg := debug.GetDebugFromContext(ctx); dbg != nil {
        dbg.AddStep("product.GetById")
    }

    return s.repository.GetById(ctx, id)
}

func (s *Service) Create(ctx context.Context, m *Product) (*Product, error) {
    return s.repository.Create(ctx, m)
}

func (s *Service) Update(ctx context.Context, m *Product) (*Product, error) {
    if err := s.repository.Update(ctx, m); err != nil {
        return nil, err
    }

    return m, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
    return s.repository.Delete(ctx, id)
}

func (s *Service) Activate(ctx context.Context, id uint) error {
    return s.repository.Activate(ctx, id)
}

func (s *Service) Deactivate(ctx context.Context, id uint) error {
    return s.repository.Deactivate(ctx, id)
}
```

---

## 2. Infrastructure — Postgres

**`internal/infrastructure/postgres/catalog/product/model.go`**

```go
package product

type product struct {
    ID         uint    `gorm:"column:id;primaryKey;autoIncrement"`
    Name       string  `gorm:"column:name"`
    Price      float64 `gorm:"column:price"`
    CategoryID uint    `gorm:"column:category_id"`
    Status     int     `gorm:"column:status"`
}

func (product) TableName() string {
    return "products"
}
```

---

**`internal/infrastructure/postgres/catalog/product/mapper.go`**

```go
package product

import domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"

func modelToEntity(m *product) *domain.Product {
    if m == nil {
        return nil
    }

    return &domain.Product{
        ID:         m.ID,
        Name:       m.Name,
        Price:      m.Price,
        CategoryID: m.CategoryID,
        Status:     m.Status,
    }
}

func entityToModel(e *domain.Product) *product {
    if e == nil {
        return nil
    }

    return &product{
        ID:         e.ID,
        Name:       e.Name,
        Price:      e.Price,
        CategoryID: e.CategoryID,
        Status:     e.Status,
    }
}
```

---

**`internal/infrastructure/postgres/catalog/product/repository.go`**

```go
package product

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/exgamer/gosdk-db-core/pkg/query/helpers"
    "github.com/exgamer/gosdk-db-core/pkg/query/pagination"
    domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"
    "gorm.io/gorm"
)

func NewPostgresRepository(client *gorm.DB) *PostgresRepository {
    return &PostgresRepository{client: client}
}

type PostgresRepository struct {
    client *gorm.DB
}

func (r *PostgresRepository) Paginated(ctx context.Context, s *domain.Search) (*pagination.Paginated[domain.Product], error) {
    helper := helpers.NewGormPaginatedHelper[domain.Product](ctx, r.client).SetPerPage(s.PerPage)

    return helper.Paginated(s.Page, func(db *gorm.DB) *gorm.DB {
        var query []string
        var args []any

        if s.ID > 0 {
            query = append(query, "products.id = ?")
            args = append(args, s.ID)
        }

        if s.Name != "" {
            query = append(query, "products.name ILIKE ?")
            args = append(args, "%"+s.Name+"%")
        }

        if s.CategoryID > 0 {
            query = append(query, "products.category_id = ?")
            args = append(args, s.CategoryID)
        }

        return db.
            WithContext(ctx).
            Select("*").
            Where(strings.Join(query, " AND "), args...).
            Order("id DESC")
    })
}

func (r *PostgresRepository) GetById(ctx context.Context, id uint) (*domain.Product, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    var m product
    result := r.client.WithContext(ctx).Where("id = ?", id).First(&m)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if result.Error != nil {
        return nil, result.Error
    }

    return modelToEntity(&m), nil
}

func (r *PostgresRepository) Create(ctx context.Context, e *domain.Product) (*domain.Product, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    m := entityToModel(e)
    if err := r.client.WithContext(ctx).Create(m).Error; err != nil {
        return nil, err
    }

    return modelToEntity(m), nil
}

func (r *PostgresRepository) Update(ctx context.Context, e *domain.Product) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return r.client.WithContext(ctx).Save(entityToModel(e)).Error
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return r.client.WithContext(ctx).Delete(product{}, "id = ?", id).Error
}

func (r *PostgresRepository) Activate(ctx context.Context, id uint) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return r.client.WithContext(ctx).Model(product{}).Where("id = ?", id).Update("status", 1).Error
}

func (r *PostgresRepository) Deactivate(ctx context.Context, id uint) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return r.client.WithContext(ctx).Model(product{}).Where("id = ?", id).Update("status", 0).Error
}
```

---

## 3. Infrastructure — Redis

**`internal/infrastructure/redis/catalog/product/repository.go`**

```go
package product

import (
    "context"
    "fmt"
    "strconv"
    "time"

    rediscore "github.com/exgamer/gosdk-redis-core/pkg/redis"
    domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"
    "github.com/redis/go-redis/v9"
)

const keyPrefix = "catalog:product:"

func NewRedisRepository(client *redis.Client) *RedisRepository {
    return &RedisRepository{client: client}
}

type RedisRepository struct {
    client *redis.Client
}

func (r *RedisRepository) Set(ctx context.Context, m *domain.Product) error {
    return rediscore.NewHelper[domain.Product](ctx, r.client).
        SetStruct(keyPrefix+strconv.FormatUint(uint64(m.ID), 10), *m, 2*time.Hour)
}

func (r *RedisRepository) GetById(ctx context.Context, id uint) (*domain.Product, error) {
    return rediscore.NewHelper[domain.Product](ctx, r.client).
        GetStruct(keyPrefix + strconv.FormatUint(uint64(id), 10))
}

func (r *RedisRepository) Invalidate(ctx context.Context, id uint) error {
    key := keyPrefix + strconv.FormatUint(uint64(id), 10)
    if err := r.client.Unlink(ctx, key).Err(); err != nil {
        return fmt.Errorf("invalidate product cache: %w", err)
    }

    return nil
}
```

---

## 4. Infrastructure — HTTP Client

**`internal/infrastructure/http/catalog/product/model.go`**

```go
package product

type productModel struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

---

**`internal/infrastructure/http/catalog/product/mapper.go`**

```go
package product

import domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"

func modelToEntity(m *productModel) *domain.Product {
    if m == nil {
        return nil
    }

    return &domain.Product{
        ID:    m.ID,
        Name:  m.Name,
        Price: m.Price,
    }
}
```

---

**`internal/infrastructure/http/catalog/product/repository.go`**

```go
package product

import (
    "context"

    "github.com/exgamer/gosdk-http-core/pkg/gin"
    "github.com/exgamer/gosdk-http-core/pkg/constants"
    "github.com/exgamer/gosdk-http-request-builder/pkg/builder"
    domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"
)

func NewHttpRepository() *HttpRepository {
    return &HttpRepository{}
}

type HttpRepository struct{}

func (r *HttpRepository) GetById(ctx context.Context, id uint) (*domain.Product, error) {
    httpInfo := gin.GetHttpInfoFromContext(ctx)

    resp, err := builder.NewGetHttpRequestBuilder[builder.Response[productModel]](
        ctx,
        "http://catalog-service/api/v1/products/1",
    ).SetRequestHeaders(map[string]string{
        constants.RequestIdHeaderName: httpInfo.RequestId,
    }).GetResult()

    if err != nil {
        return nil, err
    }

    return modelToEntity(&resp.Result.Data), nil
}
```

---

## 5. Transport — HTTP Handler

**`internal/transport/admin/http/catalog/product/request.go`**

```go
package product

import "github.com/exgamer/gosdk-http-core/pkg/validation"

type indexRequest struct {
    validation.Request
    ID         uint   `form:"id"`
    Name       string `form:"name"`
    CategoryID uint   `form:"category_id"`
    Page       uint   `form:"page"`
    PerPage    uint   `form:"per_page"`
}

type createRequest struct {
    validation.Request
    Name       string  `json:"name"        binding:"required" validate:"required"`
    Price      float64 `json:"price"       binding:"required" validate:"required,gt=0"`
    CategoryID uint    `json:"category_id" binding:"required" validate:"required"`
}

type updateRequest struct {
    validation.Request
    Name       string  `json:"name"`
    Price      float64 `json:"price"       validate:"omitempty,gt=0"`
    CategoryID uint    `json:"category_id"`
}
```

---

**`internal/transport/admin/http/catalog/product/response.go`**

```go
package product

import (
    "github.com/exgamer/gosdk-db-core/pkg/query/pagination"
    "github.com/exgamer/gosdk-http-core/pkg/structures"
)

type item struct {
    ID         uint    `json:"id"`
    Name       string  `json:"name"`
    Price      float64 `json:"price"`
    CategoryID uint    `json:"category_id"`
    Status     int     `json:"status"`
}

type itemResponse struct {
    structures.Response[item]
}

type paginatedData struct {
    Items      []*item               `json:"items"`
    Pagination pagination.Pagination `json:"pagination"`
}

type paginatedResponse struct {
    structures.Response[paginatedData]
}
```

---

**`internal/transport/admin/http/catalog/product/mapper.go`**

```go
package product

import (
    "github.com/exgamer/gosdk-db-core/pkg/query/pagination"
    domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"
)

func searchFromIndexRequest(req indexRequest) *domain.Search {
    return &domain.Search{
        ID:         req.ID,
        Name:       req.Name,
        CategoryID: req.CategoryID,
        Page:       req.Page,
        PerPage:    req.PerPage,
    }
}

func entityFromCreateRequest(req createRequest) *domain.Product {
    return &domain.Product{
        Name:       req.Name,
        Price:      req.Price,
        CategoryID: req.CategoryID,
    }
}

func entityFromUpdateRequest(id uint, req updateRequest) *domain.Product {
    return &domain.Product{
        ID:         id,
        Name:       req.Name,
        Price:      req.Price,
        CategoryID: req.CategoryID,
    }
}

func itemFromEntity(e *domain.Product) *item {
    return &item{
        ID:         e.ID,
        Name:       e.Name,
        Price:      e.Price,
        CategoryID: e.CategoryID,
        Status:     e.Status,
    }
}

func paginatedDataFromResult(p *pagination.Paginated[domain.Product]) *paginatedData {
    items := make([]*item, 0, len(p.Items))
    for _, e := range p.Items {
        items = append(items, itemFromEntity(e))
    }

    return &paginatedData{
        Items:      items,
        Pagination: *p.Pagination,
    }
}
```

---

**`internal/transport/admin/http/catalog/product/handler.go`**

```go
package product

import (
    "strconv"

    "github.com/exgamer/gosdk-http-core/pkg/response"
    "github.com/exgamer/gosdk-http-core/pkg/validators"
    domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"
    "github.com/gin-gonic/gin"
)

func NewHandler(service *domain.Service) *Handler {
    return &Handler{service: service}
}

type Handler struct {
    service *domain.Service
}

func (h *Handler) Index() gin.HandlerFunc {
    return func(c *gin.Context) {
        req := indexRequest{}
        if ok := validators.ValidateRequestQuery(c, &req); !ok {
            return
        }

        result, err := h.service.Paginated(c.Request.Context(), searchFromIndexRequest(req))
        if err != nil {
            response.InternalServerError(c, err, nil)
            return
        }

        response.Success(c, paginatedDataFromResult(result))
    }
}

func (h *Handler) View() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            response.BadRequest(c, err, nil)
            return
        }

        m, err := h.service.GetById(c.Request.Context(), uint(id))
        if err != nil {
            response.InternalServerError(c, err, nil)
            return
        }

        if m == nil {
            response.NotFound(c, nil, nil)
            return
        }

        response.Success(c, itemFromEntity(m))
    }
}

func (h *Handler) Create() gin.HandlerFunc {
    return func(c *gin.Context) {
        req := createRequest{}
        if ok := validators.ValidateRequestBody(c, &req); !ok {
            return
        }

        m, err := h.service.Create(c.Request.Context(), entityFromCreateRequest(req))
        if err != nil {
            response.InternalServerError(c, err, nil)
            return
        }

        response.SuccessCreated(c, itemFromEntity(m))
    }
}

func (h *Handler) Update() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            response.BadRequest(c, err, nil)
            return
        }

        req := updateRequest{}
        if ok := validators.ValidateRequestBody(c, &req); !ok {
            return
        }

        m, err := h.service.Update(c.Request.Context(), entityFromUpdateRequest(uint(id), req))
        if err != nil {
            response.InternalServerError(c, err, nil)
            return
        }

        response.Success(c, itemFromEntity(m))
    }
}

func (h *Handler) Delete() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            response.BadRequest(c, err, nil)
            return
        }

        if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
            response.InternalServerError(c, err, nil)
            return
        }

        response.SuccessDeleted(c, nil)
    }
}
```

---

**`internal/transport/admin/http/catalog/product/routes.go`**

```go
package product

import (
    "github.com/exgamer/gosdk-core/pkg/app"
    "github.com/exgamer/gosdk-http-core/pkg/middleware"
    httpDi "github.com/exgamer/gosdk-http-core/pkg/di"
)

func SetRoutes(a *app.App, handler *Handler) error {
    router, err := httpDi.GetRouter(a.Container)
    if err != nil {
        return err
    }

    service := router.Group("/catalog")
    {
        service.Use(middleware.RequestInfoMiddleware(a))
        service.Use(middleware.LoggerMiddleware())
        service.Use(middleware.DebugMiddleware())
        service.Use(middleware.SentryMiddleware())

        v1 := service.Group("/v1")
        {
            v1.Use(middleware.FormattedResponseMiddleware())
            v1.Use(middleware.MetricsMiddleware(a))

            v1.GET("/products", handler.Index())
            v1.GET("/product/:id", handler.View())
            v1.POST("/product", handler.Create())
            v1.PUT("/product/:id", handler.Update())
            v1.DELETE("/product/:id", handler.Delete())
        }
    }

    return nil
}
```

---

## 6. Transport — RabbitMQ Consumer

**`internal/transport/consumer/catalog/product/product_consumer.go`**

```go
package product

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/ThreeDotsLabs/watermill/message"
)

type productMessage struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

func NewConsumer() *Consumer {
    return &Consumer{}
}

type Consumer struct{}

func (c *Consumer) Consume(ctx context.Context, msg *message.Message) error {
    var payload productMessage
    if err := json.Unmarshal(msg.Payload, &payload); err != nil {
        return fmt.Errorf("unmarshal product message: %w", err)
    }

    fmt.Printf("received product: id=%d name=%s\n", payload.ID, payload.Name)

    return nil
}
```

---

**`internal/transport/consumer/catalog/product/consumer_registry.go`**

```go
package product

import "github.com/exgamer/gosdk-rabbit-core/pkg/config"

func GetConsumers(consumer *Consumer) []config.HandlerRegister {
    return []config.HandlerRegister{
        {
            Handler: consumer.Consume,
            Config: config.NewConsumerTopicDurableConfig(
                "catalog-product-consumer", // consumer name
                "catalog-exchange",          // exchange
                "catalog-product-queue",     // queue
                "catalog.product.*",         // routing key
                10,                          // workers
            ),
        },
    }
}
```

---

## 7. Bootstrap — Module

**`internal/app/bootstrap/product/repositories_factory.go`**

```go
package product

import (
    productpostgres "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/postgres/catalog/product"
    "gorm.io/gorm"
)

type repositoriesFactory struct {
    PostgresRepository *productpostgres.PostgresRepository
}

func newRepositoriesFactory(client *gorm.DB) *repositoriesFactory {
    return &repositoriesFactory{
        PostgresRepository: productpostgres.NewPostgresRepository(client),
    }
}
```

---

**`internal/app/bootstrap/product/services_factory.go`**

```go
package product

import domain "github.com/exgamer/go-sdk-rest-template/internal/domains/catalog/product"

type servicesFactory struct {
    ProductService *domain.Service
}

func newServicesFactory(repos *repositoriesFactory) *servicesFactory {
    return &servicesFactory{
        ProductService: domain.NewService(repos.PostgresRepository),
    }
}
```

---

**`internal/app/bootstrap/product/handlers_factory.go`**

```go
package product

import transport "github.com/exgamer/go-sdk-rest-template/internal/transport/admin/http/catalog/product"

type handlersFactory struct {
    ProductHandler *transport.Handler
}

func newHandlersFactory(services *servicesFactory) *handlersFactory {
    return &handlersFactory{
        ProductHandler: transport.NewHandler(services.ProductService),
    }
}
```

---

**`internal/app/bootstrap/product/consumers_factory.go`**

```go
package product

import consumer "github.com/exgamer/go-sdk-rest-template/internal/transport/consumer/catalog/product"

type consumersFactory struct {
    ProductConsumer *consumer.Consumer
}

func newConsumersFactory() *consumersFactory {
    return &consumersFactory{
        ProductConsumer: consumer.NewConsumer(),
    }
}
```

---

**`internal/app/bootstrap/product/module.go`**

```go
package product

import (
    "github.com/exgamer/gosdk-core/pkg/app"
    postgresDi "github.com/exgamer/gosdk-postgres-core/pkg/di"
    rabbitDi   "github.com/exgamer/gosdk-rabbit-core/pkg/di"
    transport  "github.com/exgamer/go-sdk-rest-template/internal/transport/admin/http/catalog/product"
    consumer   "github.com/exgamer/go-sdk-rest-template/internal/transport/consumer/catalog/product"
)

type Module struct{}

func (m *Module) Name() string {
    return "product"
}

func (m *Module) Init(a *app.App) error {
    client, err := postgresDi.GetDefaultPostgresConnection(a.Container)
    if err != nil {
        return err
    }

    repoFactory := newRepositoriesFactory(client)
    svcFactory := newServicesFactory(repoFactory)
    hdlFactory := newHandlersFactory(svcFactory)

    if err = transport.SetRoutes(a, hdlFactory.ProductHandler); err != nil {
        return err
    }

    consumersFactory := newConsumersFactory()
    consumers := consumer.GetConsumers(consumersFactory.ProductConsumer)

    reg, err := rabbitDi.GetRabbitConsumersRegistry(a.Container)
    if err != nil {
        return err
    }

    reg.RegisterMultipleHandler(consumers)

    return nil
}
```

---

## 8. Регистрация в App

**`internal/app/app.go`** — добавить модуль в `RegisterAndInitModules`:

```go
err = appInstance.RegisterAndInitModules(
    &city.Module{},
    &product.Module{},  // ← новый модуль
)
```
