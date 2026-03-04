# go-sdk-rest-template

`go-sdk-rest-template` — шаблон REST-сервиса на Go, построенный поверх `gosdk-core`, `gosdk-db-core`, `gosdk-postgres-core`, `gosdk-http-request-builder` и `gosdk-http-core`.
Предназначен для быстрого старта новых сервисов с единой архитектурой, DI, модулями и HTTP kernel.

Шаблон подходит для production-сервисов и внутренних микросервисов.

---
## Для более детализированной информации ознакомьтесь в документацей к следующим пакетам 

 - [Документация GOSDK-CORE](https://github.com/exgamer/gosdk-core)
 - [Документация GOSDK-HTTP-CORE](https://github.com/exgamer/gosdk-http-core)
 - [Документация GOSDK-DB-CORE](https://github.com/exgamer/gosdk-db-core)
 - [Документация GOSDK-POSTGRES-CORE](https://github.com/exgamer/gosdk-postgres-core)
 - [Документация GOSDK-HTTP-REQUEST-BUILDER](https://github.com/exgamer/gosdk-http-request-builder)
--- 

## 🚀 Что даёт шаблон

- 🧠 Единый `App` lifecycle
- 🌐 Готовый HTTP kernel (Gin)
- 🧩 Dependency Injection из коробки
- 🧱 Модульная архитектура (business modules)
- ⚙️ Конфигурация через env
- ♻️ Graceful shutdown
- ❗ Стандартизированная обработка ошибок
- 🧪 Удобная база для тестирования

---

## 📁 Структура проекта

[Структура проекта](STRUCTURE.md)

---

## ▶️ Запуск сервиса

### 1. Установка зависимостей

```bash
go mod tidy
```

### 2. Настройка env

```env
APP_NAME=rest-template
APP_ENV=local
APP_VERSION=1.0
TIMEZONE=
LOG_LEVEL=DEBUG

SERVER_ADDRESS=0.0.0.0:8090
SWAGGER_PREFIX=rest-template
HANDLER_TIMEOUT=30
SENTRY_DSN=

POSTGRES_DB_HOST=localhost
POSTGRES_DB_USER=admin
POSTGRES_DB_PASSWORD=admin
POSTGRES_DB_NAME=test
POSTGRES_DB_PORT=5432
POSTGRES_DB_MAX_OPEN_CONNECTIONS=100
POSTGRES_DB_MAX_IDLE_CONNECTIONS=10
POSTGRES_DB_LOG_LEVEL=info    // вывод логов SQL  запросов "info", "errors",  "warnings", если логи не нужны не указывать
```

### 3. Запуск

```bash
go run main.go
```

---

## 🧠 Архитектура

### App

`App` создаётся в `internal/app` и:
- инициализирует config
- создаёт DI контейнер
- регистрирует kernels
- запускает бизнес-модули

### Пример App

- [Пример приложения](internal/app/app.go)

### HTTP Kernel

- поднимает Gin server
- регистрирует router в DI
- управляет жизненным циклом HTTP
- [Документация GOSDK-HTTP-CORE](https://github.com/exgamer/gosdk-http-core)

### POSTGRES Kernel

- Инициализирует соединение POSTGRES
- [Документация GOSDK-POSTGRES-CORE](https://github.com/exgamer/gosdk-postgres-core)

---

## 🧱 Бизнес-модули

Каждый модуль:
- имеет статическое имя
- регистрирует свои зависимости
- добавляет HTTP-роуты

- [Пример работы с запросами](internal/infrastructure/postgres/handbook/city/repository.go)

- [Пример сервиса](internal/domains/handbook/city/service.go)

- [Пример контроллера](internal/transport/admin/http/handbook/city/handler.go)

- [Пример бизнес модуля](internal/app/bootstrap/city/module.go)

- [Пример роутов модуля](internal/transport/admin/http/handbook/city/routes.go)

---
### Режим отладки

- [Работа с Debug Collector](https://github.com/exgamer/gosdk-core/blob/main/pkg/debug/README.MD)

- [Логирование ](https://github.com/exgamer/gosdk-core/blob/main/pkg/logger/README.md)

- [Отладочная информация в ответе от API ](https://github.com/exgamer/gosdk-http-core/blob/main/DEBUG_MODE.MD)


---
### HTTP запросы

Для HTTP запросов используется пакет `gosdk-http-request-builder` (документацию смотри в README пакета)
- [Пример работы с HTTP запросами](internal/infrastructure/httpclient/handbook/city/repository.go)

---

### SWAGGER

- [Работа со SWAGGER](SWAGGER.MD)

---

## 🧩 Dependency Injection

Регистрация сервисов:

```go
di.Register(app.Container, &UserService{})
```

Использование:

```go
c, err := di.Resolve[*config.HttpConfig](a.Container)

if err != nil {
    return nil, err
}
```

---

## ❗ Ошибки

Используется единая модель ошибок из `gosdk-core` и `gosdk-http-core`:

HTTP kernel автоматически мапит ошибки в HTTP-ответы.

- [Обычные исключения](https://github.com/exgamer/gosdk-core/blob/main/pkg/exception/README.MD)
- [HTTP исключения](https://github.com/exgamer/gosdk-http-core/blob/main/pkg/exception/README.md)

---

## Логирование

- [Пример работы с логами](https://github.com/exgamer/gosdk-core/blob/v1.0.2/pkg/logger/README.md)

- [Пример работы с логами HTTP](https://github.com/exgamer/gosdk-http-core/blob/main/pkg/logger/README.md)

---

## ♻️ Graceful Shutdown 

Приложение корректно:
- ловит SIGTERM / SIGINT
- завершает HTTP сервер
- даёт время на завершение активных запросов

---

## 🧪 Рекомендации

- ✅ один бизнес-контекст = один модуль
- ✅ регистрировать роуты только в модулях
- ❌ не использовать глобальные переменные
- ❌ не создавать HTTP server вручную
- ✅ использовать DI для сервисов и репозиториев

---

## 📌 Используемые библиотеки

- `gosdk-core`
- `gosdk-http-core`
- `gosdk-db-core`
- `gosdk-postgres-core`
- `gosdk-http-request-builder`
- Gin
- Gorm
- Go 1.21+

---

## 📝 License

MIT или внутренняя лицензия компании
