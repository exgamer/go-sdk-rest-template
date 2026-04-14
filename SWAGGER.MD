## SWAGGER

- [Описать swaggerdoc в main](main.go)
- Импортировать в main.go	_ "file-uploader-service/docs"
- [Пример описания swaggerdoc в контроллере](internal/domains/handbook/modules/city/http/city/v1/handlers/city_handler.go)
- [Пример структуры для корректного отображения респонза ](internal/domains/handbook/modules/city/http/city/v1/responses/city_item_response.go)

### Генерация документации сваггера:
```go
go install github.com/swaggo/swag/cmd/swag@latest
```

```go
swag init -g main.go -o docs --parseDependency --parseInternal --parseDepth 1
```

Документация сваггера будет доступна по следующим эндпойнтам

http://0.0.0.0:8090/rest-template/api-docs/index.html

http://0.0.0.0:8090/rest-template/api-docs/doc.json