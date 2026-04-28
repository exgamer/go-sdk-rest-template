# Swagger

## Настройка

### 1. Аннотации в `main.go`

```go
// @title        My Service (API)
// @version      1.0
// @description  Описание сервиса
// @host         0.0.0.0:8090
// @BasePath     /
func main() { ... }
```

Импортировать сгенерированные docs:

```go
import _ "github.com/exgamer/go-sdk-rest-template/docs"
```

### 2. Аннотации в handler

```go
// @Summary      Получить город
// @Tags         city
// @Produce      json
// @Param        id  path  int  true  "ID города"
// @Success      200  {object}  cityResponse
// @Router       /rest-template/v1/city/{id} [get]
func (h *Handler) View() gin.HandlerFunc { ... }
```

Пример: [`internal/transport/admin/http/handbook/city/handler.go`](../internal/transport/admin/http/handbook/city/handler.go)

### 3. Response структуры

Response структуры для корректного отображения в Swagger должны эмбедить `structures.Response[T]`:

```go
type itemResponse struct {
    structures.Response[item]
}
```

---

## Генерация документации

Установить swag:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Сгенерировать docs:

```bash
swag init -g main.go -o docs --parseDependency --parseInternal --parseDepth 1
```

---

## Доступные URL

После запуска сервиса документация доступна по адресам:

```
http://0.0.0.0:8090/{SWAGGER_PREFIX}/api-docs/index.html
http://0.0.0.0:8090/{SWAGGER_PREFIX}/api-docs/doc.json
```

Значение `SWAGGER_PREFIX` задаётся в `.env`.
