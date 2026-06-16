.PHONY: run build test lint docker-up docker-down tidy

run:
	go run ./cmd/server

build:
	go build -ldflags="-s -w" -o bin/server ./cmd/server

test:
	go test ./... -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v
