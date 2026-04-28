# Правила и соглашения

## Рекомендации

- Один бизнес-контекст — один модуль
- Регистрировать маршруты только в `bootstrap/{module}/module.go`
- Не использовать глобальные переменные
- Не создавать HTTP server вручную — только через `HttpKernel`
- Использовать DI для всех зависимостей

---

## Обязательные правила

1. **Контекст** — `ctx context.Context` первым параметром во все методы
2. **Timeout** — каждый DB-запрос: `context.WithTimeout(ctx, 10*time.Second)`
3. **Not Found** — repository возвращает `nil, nil`, не ошибку
4. **Маппинг** — отдельный `mapper.go` на каждом слое, никаких прямых зависимостей структур между слоями
5. **Request DTO** — обязательно эмбедить `validation.Request`
6. **Response DTO** — обязательно эмбедить `structures.Response[T]`
7. **HTTP ответы** — только через `response.*` функции, никаких `c.JSON(...)` напрямую
8. **Ошибки** — только пробрасываются вверх; HTTP ответы пишет только transport слой

---

## Запрещённые зависимости

```
Domain → Infrastructure    ✗
Domain → Transport         ✗
Transport → Bootstrap      ✗
Infrastructure → Transport ✗
```

---

## Именование файлов

| Слой | Что | Имя файла |
|---|---|---|
| Domain | Доменная модель | `entity.go` |
| Domain | DTO поиска/фильтрации | `search.go` |
| Domain | Интерфейс репозитория | `repository.go` |
| Domain | Сервис | `service.go` |
| Infrastructure | GORM модель | `model.go` |
| Infrastructure | Реализация репозитория | `repository.go` |
| Infrastructure | Маппинг | `mapper.go` |
| Transport | Gin handler | `handler.go` |
| Transport | Регистрация маршрутов | `routes.go` |
| Transport | Request DTOs | `request.go` |
| Transport | Response DTOs | `response.go` |
| Transport | Маппинг DTO ↔ Domain | `mapper.go` |
| Transport | Consumer | `{module}_consumer.go` |
| Transport | Регистрация consumers | `consumer_registry.go` |
| Bootstrap | Точка входа модуля | `module.go` |
| Bootstrap | Фабрика | `{type}_factory.go` |
| Workflow | Оркестратор | `workflow.go` |
| Workflow | Входные/выходные данные | `dto.go` |

---

## Структура путей

| Слой | Путь |
|---|---|
| Domain | `internal/domains/{domain}/{module}/` |
| Workflow | `internal/workflow/{name}/` |
| Transport HTTP | `internal/transport/admin/http/{domain}/{module}/` |
| Transport Consumer | `internal/transport/consumer/{domain}/{module}/` |
| Infrastructure Postgres | `internal/infrastructure/postgres/{domain}/{module}/` |
| Infrastructure Redis | `internal/infrastructure/redis/{domain}/{module}/` |
| Infrastructure HTTP | `internal/infrastructure/http/{domain}/{module}/` |
| Bootstrap | `internal/app/bootstrap/{module}/` |

> Bootstrap — только уровень модуля, без домена.
> Пример: `domains/handbook/city/`, `bootstrap/city/`.

---

## Именование ключей Redis

```
{domain}:{module}:{id}

# Примеры:
catalog:product:123
handbook:city:list
```

---

## Checklist нового модуля

**Обязательно:**
- [ ] `domains/{domain}/{module}/entity.go`
- [ ] `domains/{domain}/{module}/search.go`
- [ ] `domains/{domain}/{module}/repository.go`
- [ ] `domains/{domain}/{module}/service.go`
- [ ] `infrastructure/postgres/{domain}/{module}/model.go`
- [ ] `infrastructure/postgres/{domain}/{module}/repository.go`
- [ ] `infrastructure/postgres/{domain}/{module}/mapper.go`
- [ ] `transport/admin/http/{domain}/{module}/handler.go`
- [ ] `transport/admin/http/{domain}/{module}/routes.go`
- [ ] `transport/admin/http/{domain}/{module}/request.go`
- [ ] `transport/admin/http/{domain}/{module}/response.go`
- [ ] `transport/admin/http/{domain}/{module}/mapper.go`
- [ ] `app/bootstrap/{module}/repositories_factory.go`
- [ ] `app/bootstrap/{module}/services_factory.go`
- [ ] `app/bootstrap/{module}/handlers_factory.go`
- [ ] `app/bootstrap/{module}/module.go`
- [ ] Зарегистрировать `Module` в `internal/app/app.go`

**Опционально:**
- [ ] `infrastructure/redis/{domain}/{module}/repository.go`
- [ ] `infrastructure/http/{domain}/{module}/model.go`
- [ ] `infrastructure/http/{domain}/{module}/repository.go`
- [ ] `infrastructure/http/{domain}/{module}/mapper.go`
- [ ] `transport/consumer/{domain}/{module}/{module}_consumer.go`
- [ ] `transport/consumer/{domain}/{module}/consumer_registry.go`
- [ ] `app/bootstrap/{module}/consumers_factory.go`

**Cross-domain (Workflow):**
- [ ] `workflow/{name}/workflow.go`
- [ ] `workflow/{name}/dto.go`
- [ ] `app/bootstrap/{name}/workflow_factory.go`
- [ ] Зарегистрировать workflow через DI в `module.go`
