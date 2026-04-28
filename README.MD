# go-sdk-rest-template

Шаблон REST-сервиса на Go, построенный поверх Go SDK. Предназначен для быстрого старта новых сервисов с единой архитектурой, DI, модулями и HTTP kernel.

---

## Документация

| Файл | Содержание |
|---|---|
| [ARCHITECTURE.md](documentation/ARCHITECTURE.md) | Архитектура, слои, структура проекта, жизненный цикл App |
| [NEW_DOMAIN_GUIDE.md](documentation/NEW_DOMAIN_GUIDE.md) | Пошаговое создание нового модуля с полными примерами кода |
| [SDK_REFERENCE.md](documentation/SDK_REFERENCE.md) | Справочник: DI, response, ошибки, валидация, пагинация, Redis, RabbitMQ, HTTP client, debug, middleware, конфиг |
| [CONVENTIONS.md](documentation/CONVENTIONS.md) | Правила, соглашения по именованию, checklist нового модуля |
| [SWAGGER.MD](documentation/SWAGGER.MD) | Настройка и генерация Swagger документации |

---

## Быстрый старт

### 1. Установка зависимостей

```bash
go mod tidy
```

### 2. Настройка `.env`

Скопируйте `.env.example` в `.env` и заполните значения:

```env
APP_NAME=rest-template
APP_ENV=local
APP_VERSION=1.0
LOG_LEVEL=DEBUG

SERVER_ADDRESS=0.0.0.0:8090
SWAGGER_PREFIX=rest-template
HANDLER_TIMEOUT=30

POSTGRES_DB_HOST=localhost
POSTGRES_DB_USER=admin
POSTGRES_DB_PASSWORD=admin
POSTGRES_DB_NAME=test
POSTGRES_DB_PORT=5432
POSTGRES_DB_MAX_OPEN_CONNECTIONS=100
POSTGRES_DB_MAX_IDLE_CONNECTIONS=10
POSTGRES_DB_LOG_LEVEL=info
```

### 3. Запуск

```bash
go run main.go
```

