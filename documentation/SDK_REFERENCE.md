# SDK Reference

## Содержание

1. [Dependency Injection](#1-dependency-injection)
2. [HTTP Response Helpers](#2-http-response-helpers)
3. [Обработка ошибок](#3-обработка-ошибок)
4. [Валидация запросов](#4-валидация-запросов)
5. [Пагинация](#5-пагинация)
6. [Redis кэширование](#6-redis-кэширование)
7. [RabbitMQ](#7-rabbitmq)
8. [HTTP Client Builder](#8-http-client-builder)
9. [Debug Collector](#9-debug-collector)
10. [Middleware стек](#10-middleware-стек)
11. [Конфигурация (.env)](#11-конфигурация-env)

---

## 1. Dependency Injection

DI-контейнер — тип-безопасный, основан на дженериках.

### Регистрация и получение

```go
// Зарегистрировать
di.Register(a.Container, myService)

// Получить
svc, err := di.Resolve[*MyService](a.Container)
```

### SDK DI helpers

```go
// Postgres
client, err := postgresDi.GetDefaultPostgresConnection(a.Container)

// HTTP router
router, err := httpDi.GetRouter(a.Container)

// RabbitMQ
consumersReg, err := rabbitDi.GetRabbitConsumersRegistry(a.Container)
publishersReg, err := rabbitDi.GetRabbitPublishersRegistry(a.Container)
```

---

## 2. HTTP Response Helpers

Пакет: `github.com/exgamer/gosdk-http-core/pkg/response`

Все HTTP ответы пишутся **только** через эти функции. Прямой вызов `c.JSON(...)` запрещён.

### Успешные ответы

```go
response.Success(c, data)         // 200 OK
response.SuccessCreated(c, data)  // 201 Created
response.SuccessDeleted(c, data)  // 200 OK (data можно передать nil)
```

### Ответы с ошибками

```go
response.BadRequest(c, err, nil)           // 400
response.Unauthorized(c, err, nil)         // 401
response.Forbidden(c, err, nil)            // 403
response.NotFound(c, err, nil)             // 404
response.Conflict(c, err, nil)             // 409
response.UnprocessableEntity(c, err, nil)  // 422
response.TooManyRequests(c, err, nil)      // 429
response.InternalServerError(c, err, nil)  // 500

// Универсальный — определяет статус по типу AppException автоматически
response.ErrorResponse(c, err)
```

### Формат ответов

```json
// Успех
{
  "success": true,
  "data": { ... },
  "debug": { ... }
}

// Ошибка
{
  "status": 422,
  "error": "validation",
  "message": "validation error",
  "request_id": "uuid",
  "hostname": "service-name",
  "details": { ... },
  "debug": { ... }
}
```

> `debug` присутствует только при `LOG_LEVEL=DEBUG`.

---

## 3. Обработка ошибок

Пакет: `github.com/exgamer/gosdk-core/pkg/exception`

### AppException типы

```go
// Validation → HTTP 422
exception.NewValidationException(map[string]any{"field": "message"}, false)

// Not Found → HTTP 404
exception.NewNotFoundException(errors.New("not found"), false)

// Forbidden → HTTP 403
exception.NewForbiddenException(errors.New("access denied"), false)

// Internal → HTTP 500
exception.NewAppException(err, map[string]any{"context": "value"}, true)
```

Второй булевый параметр — `trackInSentry`: отправлять ли ошибку в Sentry.

### Паттерн в handler

```go
func (h *Handler) View() gin.HandlerFunc {
    return func(c *gin.Context) {
        m, err := h.service.GetById(c.Request.Context(), id)
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
```

> Repository при записи не найдена возвращает `nil, nil` — не ошибку.

---

## 4. Валидация запросов

Пакет: `github.com/exgamer/gosdk-http-core/pkg/validators`

### Query параметры (GET)

```go
req := indexRequest{}
if ok := validators.ValidateRequestQuery(c, &req); !ok {
    return  // ответ с 422 уже отправлен автоматически
}
```

### Body (POST/PUT)

```go
req := createRequest{}
if ok := validators.ValidateRequestBody(c, &req); !ok {
    return
}
```

### Структура Request DTO

```go
type createRequest struct {
    validation.Request                               // обязательный embed
    Name  string  `json:"name"  binding:"required" validate:"required"`
    Price float64 `json:"price" binding:"required" validate:"required,gt=0"`
    Email string  `json:"email" binding:"required" validate:"required,email"`
}
```

### Кастомные правила валидации

```go
func (r createRequest) CustomValidationRules() map[string]validator.ValidationFunc {
    return map[string]validator.ValidationFunc{
        "my_rule": func(fl validator.FieldLevel) bool {
            return fl.Field().String() != "forbidden"
        },
    }
}
```

---

## 5. Пагинация

Пакет: `github.com/exgamer/gosdk-db-core/pkg/query/helpers`

### В репозитории

```go
func (r *PostgresRepository) Paginated(ctx context.Context, s *domain.Search) (*pagination.Paginated[domain.Product], error) {
    helper := helpers.NewGormPaginatedHelper[domain.Product](ctx, r.client).
        SetPerPage(s.PerPage)

    return helper.Paginated(s.Page, func(db *gorm.DB) *gorm.DB {
        return db.
            WithContext(ctx).
            Select("*").
            Where("status = ?", 1).
            Order("id DESC")
    })
}
```

### Response структура

```go
type paginatedData struct {
    Items      []*item               `json:"items"`
    Pagination pagination.Pagination `json:"pagination"`
}
```

### JSON ответ

```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "total": 100,
      "last_page": 10,
      "current_page": 1,
      "per_page": 10,
      "prev_page": 0,
      "next_page": 2
    }
  }
}
```

---

## 6. Redis кэширование

Пакет: `github.com/exgamer/gosdk-redis-core/pkg/redis`

```go
helper := rediscore.NewHelper[domain.Product](ctx, r.client)

// Struct
helper.SetStruct("key", product, 2*time.Hour)
product, err := helper.GetStruct("key")  // nil, nil если не найдено

// Массив
helper.SetArray("key", products, 1*time.Hour)
products, err := helper.GetArray("key")

// Строка
helper.SetString("key", "value", 30*time.Minute)
value, err := helper.GetString("key")
```

Схема ключей: `{domain}:{module}:{id}`, например: `catalog:product:123`.

---

## 7. RabbitMQ

### Publisher

```go
publishersReg, err := rabbitDi.GetRabbitPublishersRegistry(a.Container)
publisher := publishersReg.GetPublisher()

msg := message.NewMessage(watermill.NewUUID(), payload)
publisher.Publish("exchange-name", msg)
```

### Consumer

```go
func (c *Consumer) Consume(ctx context.Context, msg *message.Message) error {
    // nil = ACK (сообщение обработано)
    // error = NACK (сообщение будет повторено)
    return nil
}
```

### Конфигурация очереди

```go
config.NewConsumerTopicDurableConfig(
    "consumer-name",  // уникальное имя
    "exchange-name",  // exchange
    "queue-name",     // очередь
    "routing.key.*",  // routing key (поддерживает wildcards)
    10,               // кол-во параллельных воркеров
)
```

---

## 8. HTTP Client Builder

Пакет: `github.com/exgamer/gosdk-http-request-builder/pkg/builder`

### Методы

```go
// GET
resp, err := builder.NewGetHttpRequestBuilder[ResponseType](ctx, url).
    SetRequestHeaders(headers).
    SetQueryParams(map[string]string{"page": "1"}).
    GetResult()

// POST
resp, err := builder.NewPostHttpRequestBuilder[ResponseType](ctx, url).
    SetJSONBody(requestBody).
    SetRequestTimeout(5 * time.Second).
    GetResult()

// PUT
resp, err := builder.NewPutHttpRequestBuilder[ResponseType](ctx, url).
    SetJSONBody(requestBody).
    GetResult()

// DELETE
resp, err := builder.NewDeleteHttpRequestBuilder[ResponseType](ctx, url).
    GetResult()
```

### Проверка статуса

```go
resp.IsSuccess()      // 2xx
resp.IsClientError()  // 4xx
resp.IsServerError()  // 5xx
```

### Передача Request ID

```go
httpInfo := gin.GetHttpInfoFromContext(ctx)

builder.NewGetHttpRequestBuilder[ResponseType](ctx, url).
    SetRequestHeaders(map[string]string{
        constants.RequestIdHeaderName: httpInfo.RequestId,
    })
```

---

## 9. Debug Collector

Пакет: `github.com/exgamer/gosdk-core/pkg/debug`

Собирает шаги выполнения запроса и добавляет их в ответ при `LOG_LEVEL=DEBUG`. Инициализируется через `DebugMiddleware()`. SDK автоматически добавляет SQL-запросы и HTTP-вызовы.

```go
if dbg := debug.GetDebugFromContext(ctx); dbg != nil {
    dbg.AddStep("product.GetById: cache miss, fetching from db")
}
```

---

## 10. Middleware стек

Пакет: `github.com/exgamer/gosdk-http-core/pkg/middleware`

Порядок регистрации важен:

```go
service := router.Group("/prefix")

// Уровень сервиса
service.Use(middleware.RequestInfoMiddleware(a))  // request ID, IP, method, path
service.Use(middleware.LoggerMiddleware())         // структурированное логирование
service.Use(middleware.DebugMiddleware())          // debug collector
service.Use(middleware.SentryMiddleware())         // трекинг ошибок

v1 := service.Group("/v1")

// Уровень версии API
v1.Use(middleware.FormattedResponseMiddleware())   // обёртка ответов
v1.Use(middleware.MetricsMiddleware(a))            // Prometheus метрики
```

| Middleware | Что делает |
|---|---|
| `RequestInfoMiddleware` | Добавляет request ID, IP, method, path в контекст |
| `LoggerMiddleware` | Логирует каждый запрос (метод, путь, статус, время) |
| `DebugMiddleware` | Создаёт DebugCollector в контексте |
| `SentryMiddleware` | Перехватывает паники и ошибки, отправляет в Sentry |
| `FormattedResponseMiddleware` | Оборачивает ответы в стандартный формат |
| `MetricsMiddleware` | Собирает Prometheus метрики (latency, status codes) |

---

## 11. Конфигурация (.env)

```env
# Приложение
APP_NAME=my-service
APP_ENV=local          # local | staging | production
APP_VERSION=1.0
TIMEZONE=UTC
LOG_LEVEL=DEBUG        # DEBUG | INFO | WARN | ERROR

# HTTP сервер
SERVER_ADDRESS=0.0.0.0:8090
SWAGGER_PREFIX=my-service
HANDLER_TIMEOUT=30

# Sentry
SENTRY_DSN=

# Postgres
POSTGRES_DB_HOST=localhost
POSTGRES_DB_USER=admin
POSTGRES_DB_PASSWORD=admin
POSTGRES_DB_NAME=mydb
POSTGRES_DB_PORT=5432
POSTGRES_DB_MAX_OPEN_CONNECTIONS=100
POSTGRES_DB_MAX_IDLE_CONNECTIONS=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# RabbitMQ
RABBIT_HOST=localhost
RABBIT_PORT=5672
RABBIT_USER=guest
RABBIT_PASSWORD=guest
RABBIT_VHOST=/
```
