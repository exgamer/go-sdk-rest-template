REST Application Example with using 
https://github.com/exgamer/go-rest-sdk

1 create table test_table.sql
2 create .env file from .env.example
3 go mod tidy
4 install swagger: go install github.com/swaggo/swag/cmd/swag@latest
5 Generate swagger docs: swag init
4 go run main.go
5 Test endpoint: curl --location 'http://0.0.0.0:8090/example-go/v1/test'
6 Metrics endpoint: curl --location 'http://0.0.0.0:8090/metrics'
7 Healthcheck readiness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/readiness'
8 Healthcheck liveness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/liveness'
9 Swagger docs curl --location 'http://localhost:8090/swagger/index.html'


