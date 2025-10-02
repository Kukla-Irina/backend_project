# Lists API (CRUD, v1)

Простой REST API для работы со списками.  
Поддерживает базовые операции: создание, получение, обновление и удаление списков.  

- Go 1.25.1
- REST JSON, `/api/v1`
- In-memory storage

## Запуск
```bash
go mod tidy
go run ./cmd/todo-api
# сервер будет доступен по адресу http://localhost:8080
```

## Примеры запросов и ответов
```bash
# health-check
curl -s http://localhost:8080/health
# Ответ:
# OK
```

### Создать список
```bash
curl -s -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"Дом"}'
# Ответ:
# {
#   "id": "e7a1c3fa-4bdf-4c67-b1b8-7c14c1d3f841",
#   "title": "Дом",
#   "created_at": "2025-10-01T12:34:56Z"
# }
```

### Получить все списки (с пагинацией)
```bash
curl -s "http://localhost:8080/api/v1/lists?limit=10&offset=0" -i
# Ответ (заголовок X-Total-Count: 1):
# [
#   {"id":"...","title":"Дом","created_at":"..."}
# ]
```

### Получить список по ID
```bash
curl -s http://localhost:8080/api/v1/lists/<uuid>
# Ответ:
# {
#   "id": "<uuid>",
#   "title": "Дом",
#   "created_at": "2025-10-01T12:34:56Z"
# }
```

### Обновить title
```bash
curl -s -X PATCH http://localhost:8080/api/v1/lists/<uuid> \
  -H "Content-Type: application/json" \
  -d '{"title":"Домашние дела"}'
# Ответ:
# {
#   "id": "<uuid>",
#   "title": "Домашние дела",
#   "created_at": "2025-10-01T12:34:56Z"
# }
```

### Удалить список
```bash
curl -s -X DELETE http://localhost:8080/api/v1/lists/<uuid> -i
# Ответ:
# HTTP/1.1 204 No Content
```
