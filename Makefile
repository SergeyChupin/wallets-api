.PHONY: build
build:
	go build -v ./cmd/httpserver

.PHONY: test
test:
	go test -v -cover -race -timeout 30s ./internal/...

.PHONY: compose-up
compose-up:
	docker-compose up --build -d

compose-logs:
	docker-compose logs -f

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: check-swagger
check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

.PHONY: swagger
swagger: check-swagger
	swagger generate spec -o ./docs/api.yaml --scan-models

.DEFAULT_GOAL := build
