.PHONY: build run swaggui install-tools ogen mockery genall testcover testunit

build:
	go build -o marketplace-parser-service ./cmd/marketplace-parser-service

run: build
	SERVER_HTTP_ADDR="localhost:8080" CONFIG_PATH=./configs/config.yaml ./marketplace-parser-service

swaggui:
	docker run --rm -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v ./api/v1/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui

install-tools:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/vektra/mockery/v3@latest

ogen:
	ogen --target ./internal/transport/http/httpgen --package httpgen --clean ./api/v1/openapi.yaml

mockery:
	mockery --config .mockery.yaml --log-level=debug

genall: ogen mockery

testcover:
	go test -coverprofile=coverage.out ./...

testunit:
	go test ./internal/usecase/
	go test ./internal/transport/http
	go test ./internal/adapters/browser/chromium