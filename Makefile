.DEFAULT_GOAL := run

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run ./...
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

run: vet
	go run ./cmd/main.go
.PHONY:run

tidy:
	go mod tidy
.PHONY:tidy

tests:
	go test ./...
.PHONY:tests

up:
	docker-compose -f ./docker-compose.yml up -d --force-recreate --build
.PHONY:up

down:
	docker-compose -f ./docker-compose.yml down
.PHONY:down