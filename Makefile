.PHONY: run build tidy

run:
	go run ./cmd/todo-api

build:
	go build -o bin/todo-api ./cmd/todo-api

tidy:
	go mod tidy
