# Backend Architecture Regulation (Go Example)

## Цель документа

Этот документ описывает архитектурные правила backend‑сервисов.

Регламент помогает:

-   поддерживать единообразную структуру сервисов
-   разделять ответственность между слоями
-   уменьшать связанность компонентов
-   упрощать тестирование и развитие системы

Архитектура основана на принципах:

-   Clean Architecture
-   Dependency Inversion
-   Separation of Concerns

------------------------------------------------------------------------

# Основные принципы

## 1. Разделение слоев

Сервис должен быть разделён на логические слои.

Application (Bootstrap)\
Transport\
Domain\
Infrastructure

Каждый слой имеет чёткую ответственность.

------------------------------------------------------------------------

# Направление зависимостей

Зависимости должны направляться к бизнес‑логике.

Transport\
↓\
Domain\
↑\
Infrastructure\
↑\
Application (Bootstrap)

Это означает:

-   внешние слои могут зависеть от внутренних
-   внутренние слои не должны зависеть от внешних

------------------------------------------------------------------------

# Пример структуры проекта (Go)

    internal/

      app/
        bootstrap/
          order/
            module.go

      contexts/
        admin/
          http/
            order/
              handler.go
              routes.go

      domains/
        order/
          order.go
          repository.go
          service.go

      infrastructure/
        postgres/
          order/
            repository.go

------------------------------------------------------------------------

# Domain Layer

## Назначение

Domain слой содержит бизнес‑логику системы.

Он не должен зависеть от инфраструктуры или транспорта.

------------------------------------------------------------------------

## Пример доменной модели

``` go
package order

type Order struct {
    ID    int
    Price int
}
```

------------------------------------------------------------------------

## Пример интерфейса репозитория

``` go
package order

type Repository interface {
    Save(order Order) error
    FindByID(id int) (*Order, error)
}
```

------------------------------------------------------------------------

## Пример доменного сервиса

``` go
package order

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(order Order) error {
    return s.repo.Save(order)
}
```

Domain не знает:

-   какая база данных используется
-   какой ORM
-   какой транспорт

------------------------------------------------------------------------

# Infrastructure Layer

## Назначение

Infrastructure реализует интеграции с внешними системами.

Примеры:

-   базы данных
-   внешние API
-   брокеры сообщений
-   кэш

------------------------------------------------------------------------

## Пример реализации репозитория

``` go
package order

import (
    "database/sql"
    domain "example/internal/domains/order"
)

type PostgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(order domain.Order) error {
    _, err := r.db.Exec(
        "INSERT INTO orders(id, price) VALUES ($1,$2)",
        order.ID,
        order.Price,
    )
    return err
}

func (r *PostgresRepository) FindByID(id int) (*domain.Order, error) {
    row := r.db.QueryRow("SELECT id, price FROM orders WHERE id=$1", id)

    order := domain.Order{}
    err := row.Scan(&order.ID, &order.Price)

    return &order, err
}
```

Infrastructure зависит от Domain.

------------------------------------------------------------------------

# Transport Layer

## Назначение

Transport отвечает за взаимодействие с внешним миром.

Примеры:

-   HTTP
-   gRPC
-   Message Consumers
-   CLI

------------------------------------------------------------------------

## Пример HTTP handler

``` go
package order

import (
    "encoding/json"
    "net/http"
    domain "example/internal/domains/order"
)

type Handler struct {
    service *domain.Service
}

func NewHandler(service *domain.Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

    var req struct {
        ID    int `json:"id"`
        Price int `json:"price"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    order := domain.Order{
        ID: req.ID,
        Price: req.Price,
    }

    err := h.service.Create(order)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
```

Handler:

-   принимает HTTP запрос
-   конвертирует DTO → Domain
-   вызывает сервис

------------------------------------------------------------------------

# Application / Bootstrap Layer

## Назначение

Bootstrap слой отвечает за сборку приложения.

Именно здесь:

-   создаются подключения
-   выбираются реализации
-   собираются зависимости
-   регистрируются маршруты

------------------------------------------------------------------------

## Пример bootstrap

``` go
package bootstrap

import (
    "database/sql"
    "net/http"

    domain "example/internal/domains/order"
    postgres "example/internal/infrastructure/postgres/order"
    httpOrder "example/internal/contexts/admin/http/order"
)

func Start() {

    db, _ := sql.Open("postgres", "...")

    repo := postgres.NewPostgresRepository(db)

    service := domain.NewService(repo)

    handler := httpOrder.NewHandler(service)

    http.HandleFunc("/orders", handler.Create)

    http.ListenAndServe(":8080", nil)
}
```

Bootstrap соединяет:

Infrastructure + Domain + Transport

------------------------------------------------------------------------

# Dependency Rules

Допустимые зависимости:

Transport → Domain\
Infrastructure → Domain\
Bootstrap → Transport\
Bootstrap → Infrastructure\
Bootstrap → Domain

Запрещённые зависимости:

Domain → Infrastructure\
Domain → Transport\
Transport → Bootstrap\
Infrastructure → Transport

------------------------------------------------------------------------

# Dependency Injection

Создание зависимостей происходит в bootstrap слое.

Цепочка зависимостей:

database\
↓\
repository\
↓\
service\
↓\
handler

------------------------------------------------------------------------

# DTO и Domain модели

Transport слой использует DTO.

CreateOrderRequest\
CreateOrderResponse

Domain использует доменные модели.

Transport отвечает за преобразование:

RequestDTO → DomainModel\
DomainModel → ResponseDTO

------------------------------------------------------------------------

# Тестирование

Domain слой должен легко тестироваться без инфраструктуры.

Unit‑тесты должны покрывать:

-   сервисы
-   бизнес‑правила
-   доменные сценарии

Infrastructure тестируется отдельно.
