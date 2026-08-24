.PHONY: run build test migrate-up migrate-down tidy

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -count=1

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "${DATABASE_URL}" down 1
