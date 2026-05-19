# Архитектура сервиса

## Содержание

1. [Обзор](#1-обзор)
2. [Ответственность слоёв](#2-ответственность-слоёв)
3. [SDK пакеты](#3-sdk-пакеты)
4. [Структура проекта](#4-структура-проекта)
5. [App и Kernels](#5-app-и-kernels)
6. [Cross-domain взаимодействие](#6-cross-domain-взаимодействие)

---

## 1. Обзор

Сервис строится по принципам **Clean Architecture**. Зависимости направлены внутрь — к бизнес-логике.

```
┌─────────────────────────────────────────┐
│         Application / Bootstrap         │  ← собирает всё вместе
├─────────────────────────────────────────┤
│              Transport Layer            │  ← HTTP handlers, consumers
├─────────────────────────────────────────┤
│            Workflow Layer (*)           │  ← cross-domain оркестрация
├─────────────────────────────────────────┤
│               Domain Layer              │  ← бизнес-логика, интерфейсы
├─────────────────────────────────────────┤
│           Infrastructure Layer          │  ← Postgres, Redis, HTTP client
└─────────────────────────────────────────┘

(*) опциональный слой — используется только при cross-domain сценариях
```

**Направление зависимостей:**

```
Transport      ──→  Workflow / Domain
Workflow       ──→  Domain
Infrastructure ──→  Domain
Bootstrap      ──→  Transport + Workflow + Infrastructure + Domain
```

Domain не знает ни про базу данных, ни про HTTP — только про бизнес-правила.

**Запрещённые зависимости:**

```
Domain     → Infrastructure    ✗
Domain     → Transport         ✗
Domain     → Workflow          ✗
Workflow   → Transport         ✗
Transport  → Bootstrap         ✗
Infrastructure → Transport     ✗
```

---

## 2. Ответственность слоёв

Иллюстративные примеры структуры каждого слоя: → [`examples/layers.go`](examples/layers.go)

### Domain Layer

Содержит бизнес-логику. Не зависит ни от базы данных, ни от HTTP, ни от любой внешней системы.

Что входит:
- **Entity** — доменная модель (чистые Go-структуры без тегов ORM/JSON)
- **Repository interface** — контракт доступа к данным, определённый доменом
- **Service** — оркестрирует бизнес-операции через интерфейс репозитория
- **Search / DTO** — объекты для передачи параметров внутри домена

### Workflow Layer *(опциональный)*

Оркестрирует несколько доменных сервисов для выполнения бизнес-процесса, который не принадлежит ни одному домену.

Что входит:
- **Workflow** — последовательность вызовов доменных сервисов
- **DTO** — входные и выходные данные процесса

Правила:
- Зависит от Domain (вызывает доменные сервисы через интерфейсы)
- Не знает про HTTP, Gin, очереди
- Может содержать бизнес-логику и репозиторий если они принадлежат процессу, а не конкретному домену

Когда создавать Workflow — см. [раздел 6](#6-cross-domain-взаимодействие).

### Infrastructure Layer

Реализует интеграции с внешними системами: базами данных, кешем, внешними API.

Что входит:
- **Postgres repository** — реализация доменного интерфейса через GORM
- **Redis repository** — кеширование через gosdk-redis-core
- **HTTP repository** — вызовы внешних сервисов через gosdk-http-request-builder
- **Mapper** — преобразование между инфраструктурной моделью и доменной сущностью

Infrastructure зависит от Domain, но Domain не знает про Infrastructure.

### Transport Layer

Отвечает за взаимодействие с внешним миром: HTTP, RabbitMQ.

Что входит:
- **Handler** — принимает HTTP запрос, вызывает сервис, возвращает ответ
- **Routes** — регистрирует маршруты и middleware стек
- **Request / Response DTOs** — объекты для сериализации входных и выходных данных
- **Mapper** — преобразование Request DTO → Domain entity и Domain entity → Response DTO
- **Consumer** — обработчик сообщений из RabbitMQ

```
Handler:
  1. Валидировать входные данные
  2. Преобразовать Request DTO → Domain entity (mapper)
  3. Вызвать сервис
  4. Преобразовать Domain entity → Response DTO (mapper)
  5. Вернуть HTTP ответ через response.*
```

### Bootstrap / Application Layer

Собирает приложение: создаёт зависимости, подключает слои друг к другу.

Что входит:
- **App** — инициализирует Kernels и Modules
- **Module** — точка входа модуля, создаёт factory цепочку, регистрирует маршруты и consumers
- **Factories** — создают конкретные реализации репозиториев, сервисов, handlers

```
Цепочка зависимостей:
  DB Client → Repository (infra) → Service (domain) → Handler (transport) → Routes
```

---

## 3. SDK пакеты


| Пакет | Импорт | Назначение |
|---|---|---|
| `gosdk-core` | `github.com/exgamer/gosdk-core` | App, DI, Kernel/Module интерфейсы, AppException, DebugCollector |
| `gosdk-http-core` | `github.com/exgamer/gosdk-http-core` | HttpKernel, middleware, response helpers, validators |
| `gosdk-postgres-core` | `github.com/exgamer/gosdk-postgres-core` | PostgresKernel, DI helpers для GORM |
| `gosdk-db-core` | `github.com/exgamer/gosdk-db-core` | GormPaginatedHelper, Paginated[T], Pagination |
| `gosdk-redis-core` | `github.com/exgamer/gosdk-redis-core` | Redis helper (SetStruct/GetStruct/SetArray) |
| `gosdk-rabbit-core` | `github.com/exgamer/gosdk-rabbit-core` | RabbitKernel, ConsumersRegistry, PublisherRegistry |
| `gosdk-http-request-builder` | `github.com/exgamer/gosdk-http-request-builder` | Type-safe HTTP client builder |

**Репозитории:**
- [gosdk-core](https://github.com/exgamer/gosdk-core)
- [gosdk-http-core](https://github.com/exgamer/gosdk-http-core)
- [gosdk-db-core](https://github.com/exgamer/gosdk-db-core)
- [gosdk-postgres-core](https://github.com/exgamer/gosdk-postgres-core)
- [gosdk-http-request-builder](https://github.com/exgamer/gosdk-http-request-builder)
- [gosdk-redis-core](https://github.com/exgamer/gosdk-redis-core)
- [gosdk-rabbit-core](https://github.com/exgamer/gosdk-rabbit-core)

---

## 4. Структура проекта

```
internal/
├── app/
│   ├── app.go                              ← инициализация App, регистрация Kernels и Modules
│   └── bootstrap/
│       └── {module}/                       ← например: city, product (без домена)
│       └── {domain}/{module}/              ← например: handbook/city (с доменом, если нужна группировка)
│           ├── module.go                   ← точка входа модуля
│           ├── repositories_factory.go     ← создание репозиториев
│           ├── services_factory.go         ← создание сервисов
│           ├── handlers_factory.go         ← создание handlers
│           └── consumers_factory.go        ← создание consumers (опционально)
│
├── domains/
│   └── {domain}/                           ← например: handbook, order, payment
│       └── {module}/                       ← например: city, product, invoice
│           ├── entity.go                   ← доменная модель
│           ├── search.go                   ← DTO поиска/фильтрации
│           ├── repository.go               ← интерфейс репозитория
│           └── service.go                  ← бизнес-логика
│
├── workflow/                               ← cross-domain оркестрация (опционально)
│   └── {name}/                             ← например: checkout, registration
│       ├── workflow.go                     ← оркестрирует несколько доменных сервисов
│       └── dto.go                          ← входные/выходные данные workflow
│
├── transport/
│   ├── admin/
│   │   └── http/
│   │       └── {domain}/                   ← например: handbook
│   │           └── {module}/               ← например: city
│   │               ├── handler.go          ← gin handlers
│   │               ├── routes.go           ← регистрация маршрутов
│   │               ├── request.go          ← Request DTOs
│   │               ├── response.go         ← Response DTOs
│   │               └── mapper.go           ← DTO ↔ Domain маппинг
│   └── base/
│       └── rabbit/
│           └── {domain}/                       ← например: handbook
│               └── {module}/                   ← например: city
│                  ├── consumer_registry.go    ← регистрация consumers
│                  └── {module}_consumer.go    ← логика обработки сообщений
│
└── infrastructure/
    ├── postgres/
    │   └── {domain}/                       ← например: handbook
    │       └── {module}/                   ← например: city
    │           ├── model.go                ← GORM модель
    │           ├── repository.go           ← реализация Repository
    │           └── mapper.go               ← Model ↔ Domain маппинг
    ├── redis/
    │   └── {domain}/
    │       └── {module}/
    │           └── repository.go           ← кэш операции
    └── http/
        └── {domain}/
            └── {module}/
                ├── model.go                ← модель HTTP ответа
                ├── repository.go           ← HTTP клиент
                └── mapper.go               ← Model ↔ Domain маппинг
```

**Ключевые правила структуры:**
- `bootstrap/` — по умолчанию `{module}/`; домен добавляется если нужна группировка: `{domain}/{module}/`
- Все остальные слои всегда: `{domain}/{module}/`
- Примеры: `bootstrap/city/` или `bootstrap/handbook/city/`, `domains/handbook/city/`, `transport/admin/http/handbook/city/`

---

## 5. App и Kernels

### Жизненный цикл приложения

```
main()
  └── NewApp()
        ├── RegisterAndInitKernels(...)
        │     ├── PostgresKernel.Init()  → подключение к БД, регистрация в DI
        │     ├── HttpKernel.Init()      → создание Gin роутера, регистрация в DI
        │     └── RabbitKernel.Init()    → AMQP соединение, регистрация registry в DI
        │
        └── RegisterAndInitModules(...)
              └── Module.Init()
                    ├── Создание factory цепочки
                    ├── Регистрация HTTP маршрутов
                    └── Регистрация RabbitMQ consumers
  └── RunAll()
        ├── HttpKernel.Start()    → запуск HTTP сервера в горутине
        └── RabbitKernel.Start()  → запуск consumer в горутине
  └── WaitForShutdown()
        └── Ожидание SIGINT/SIGTERM → graceful shutdown всех Kernels
```

### Инициализация App (`internal/app/app.go`)

Примеры: регистрация Kernels, варианты RabbitKernel, шаблон Module interface:
→ [`examples/app_init.go`](examples/app_init.go)

---

## 6. Cross-domain взаимодействие

Когда бизнес-логика затрагивает несколько доменов, есть два подхода. Выбор зависит от сложности сценария.

---

### Случай 1 — прямая зависимость между сервисами

**Когда использовать:** один сервис нужен другому для простой валидации или получения данных. Связь очевидна и стабильна.

**Как:** инжектировать доменный сервис (или его интерфейс) через конструктор.

**Пример:** `OrderService` проверяет существование города перед созданием заказа.

Полный пример со всеми шагами (service.go, city/module.go, order/module.go, services_factory.go, порядок в app.go):
→ [`examples/cross_domain_direct.go`](examples/cross_domain_direct.go)

**Ключевые правила:**
- Интерфейс `CityServiceInterface` объявляется в домене-потребителе (`order`), не в провайдере
- `city.Service` реализует контракт неявно — никаких изменений в `handbook/city` не требуется
- `city.Module` регистрирует сервис в DI; `order.Module` получает его — поэтому city должен быть в `app.go` первым
- `OrderService` зависит от интерфейса, а не от `*city.Service` напрямую — домен `order` не импортирует конкретную реализацию

---

### Случай 2 — Workflow (оркестрация нескольких доменов)

**Когда использовать:** бизнес-процесс не принадлежит ни одному домену. Несколько доменов участвуют как равноправные участники операции.

**Примеры:** оформление заказа, регистрация пользователя, обработка платежа.

**Структура:**
```
internal/
└── workflow/
    └── {name}/
        ├── workflow.go   ← оркестрирует доменные сервисы
        └── dto.go        ← входные/выходные данные workflow
```

**Пример:** `CheckoutWorkflow` — оформление заказа.

Полный пример: `dto.go`, `workflow.go`, `workflow_factory.go`, вызов из handler:
→ [`examples/cross_domain_workflow.go`](examples/cross_domain_workflow.go)

---

### Когда что использовать

| Ситуация | Подход |
|---|---|
| Сервису A нужны данные из домена B | Случай 1 — прямая зависимость через интерфейс |
| Валидация по данным другого домена | Случай 1 — прямая зависимость через интерфейс |
| Операция охватывает 2+ доменов как равноправных участников | Случай 2 — Workflow |
| Бизнес-процесс из нескольких шагов с частичной компенсацией | Случай 2 — Workflow |
| Нужно переиспользовать процесс из разных точек входа (HTTP + Consumer) | Случай 2 — Workflow |
