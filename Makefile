.PHONY: run build test lint up down

run:
	go run ./cmd/server/...

build:
	go build -o bin/server ./cmd/server/...

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

up:
	docker compose up -d

down:
	docker compose down -v