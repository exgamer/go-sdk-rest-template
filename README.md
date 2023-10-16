REST Application Example with using https://github.com/exgamer/go-rest-sdk

1 [Set up environment](manual/ENVIRONMENT.MD)

2 Create table example/test_table.sql

3 Create .env file from .env.example

4 go mod tidy

5 Install swagger: go install github.com/swaggo/swag/cmd/swag@latest

6 Generate swagger docs: swag init

7 go run main.go

8 Test endpoint: curl --location 'http://0.0.0.0:8090/example-go/v1/test'

9 Metrics endpoint: curl --location 'http://0.0.0.0:8090/metrics'

10 Healthcheck readiness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/readiness'

11 Healthcheck liveness endpoint: curl --location 'http://0.0.0.0:8090/healthcheck/liveness'

12 Swagger docs curl --location 'http://localhost:8090/example-go/api-docs/index.html'
13 Swagger docs curl --location 'http://localhost:8090/example-go/api-docs/doc.json'

14 [Create docker image for your rest application](manual/HOW_TO_CREATE_DOCKER_FOR_GO_SERVICE.MD)

