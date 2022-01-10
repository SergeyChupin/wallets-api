.PHONY: build
build:
	go build -v ./cmd/httpserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: setup
setup:
	docker-compose up

.PHONY: check-swagger
check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

.PHONY: swagger
swagger: check-swagger
	swagger generate spec -o ./docs/api.yaml --scan-models

.DEFAULT_GOAL := build
