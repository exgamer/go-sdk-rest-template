## SWAGGER

- [Описать swaggerdoc в main](main.go)
- [Пример описания swaggerdoc в контроллере](internal/domains/handbook/modules/city/http/city/v1/handlers/city_handler.go)
- [Пример структуры для корректного отображения респонза ](internal/domains/handbook/modules/city/http/city/v1/responses/city_item_response.go)

### Генерация документации сваггера:
```go
go install github.com/swaggo/swag/cmd/swag@latest
```

```go
swag init -g main.go -o docs --parseDependency --parseInternal --parseDepth 1
```