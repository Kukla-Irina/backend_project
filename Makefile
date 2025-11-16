.PHONY: run build tidy migrate-up migrate-down migrate-create migrate-version

# --- Параметры окружения ---
DB_URL = postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable

# --- Go targets ---
run:
		go run ./cmd/todo-api

build:
		go build -o bin/todo-api ./cmd/todo-api

tidy:
		go mod tidy

# --- Миграции (cli migrate) ---
# Применить все миграции вверх
migrate-up:
		migrate -path migrations -database "$(DB_URL)" up

# Откатить одну миграцию
migrate-down:
		migrate -path migrations -database "$(DB_URL)" down 1

# Создать новую миграцию: make migrate-create NAME=add_items_table
migrate-create:
		migrate create -ext sql -dir migrations -seq $(NAME)

# Посмотреть текущую версию
migrate-version:
		migrate -path migrations -database "$(DB_URL)" version


