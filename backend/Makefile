.PHONY: build run test test-unit test-integration test-api lint tidy migrate-up migrate-down clean

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

test: test-unit test-integration test-api

test-unit:
	go test ./internal/... -count=1 -short

test-integration:
	go test ./internal/... -count=1 -run Integration -db-host=localhost -db-port=5433 -db-user=postgres -db-password=password -db-name=owndangan_test

test-api:
	go test ./internal/... -count=1 -run APITest -db-host=localhost -db-port=5433 -db-user=postgres -db-password=password -db-name=owndangan_test

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrate -direction=up

migrate-down:
	go run ./cmd/migrate -direction=down

clean:
	rm -rf bin/
