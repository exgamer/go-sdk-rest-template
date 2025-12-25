## SWAGGER

- [Описать swaggerdoc в main](main.go)
- [Пример описания wsaggerdoc в контроллере](internal/domains/handbook/modules/city/http/city/handlers/city_handler.go)
- [Пример структуры для корректного отображения респонза ](internal/domains/handbook/modules/city/http/city/responses/one_item_response.go)

### Генерация документации сваггера:

```go
swag init -g main.go -o docs --parseDependency --parseInternal --parseDepth 1
```