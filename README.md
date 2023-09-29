REST Application Example with using https://github.com/exgamer/go-rest-sdk

1 [set up environment](manual/ENVIRONMENT.md)

2 create table example/test_table.sql

3 create .env file from .env.example

4 go mod tidy

5 install swagger: go install github.com/swaggo/swag/cmd/swag@latest

6 Generate swagger docs: swag init

7 go run main.go

8 Test endpoint: curl --location 'http://0.0.0.0:8090/example-go/v1/test'

9 Metrics endpoint: curl --location 'http://0.0.0.0:8090/metrics'

10 Healthcheck readiness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/readiness'

11 Healthcheck liveness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/liveness'

12 Swagger docs curl --location 'http://localhost:8090/swagger/index.html'


