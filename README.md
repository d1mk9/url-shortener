# URL Shortener

Минимальный сервис для сокращения URL и управления короткими ссылками.

## Возможности
- Создание короткой ссылки (`POST /shorten`)
- Переход по короткой ссылке (`GET /{short_id}`)
- Удаление ссылки (`DELETE /{short_id}`)
- Health-check (`GET /healthz`)

## Установка и запуск

### Требования
- Go 1.22+
- PostgreSQL

### Настройки
Конфигурация задаётся через YAML‑файл (`CONFIG_FILE`) и переменные окружения:

```yaml
postgres_host: "localhost"
postgres_port: "5432"
postgres_db: "urlshortener"
base_url: "http://urlshort.xy"
```

А также переменные окружения:

```bash
export POSTGRES_USER=<username>
export POSTGRES_PASSWORD=<password>
```

### Makefile команды
В проекте есть удобный `Makefile`:

```makefile
.PHONY: run migrate-up migrate-down tidy

run:
	go run ./cmd/app serve

migrate-up:
	go run ./cmd/app migrate up

migrate-down:
	go run ./cmd/app migrate down

tidy:
	go mod tidy
```

## REST API

### Создать короткую ссылку
```bash
curl -i -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"OriginalURL":"https://example.com","MaxVisits":3,"Expiry":"24h"}'
```

### Перейти по короткой ссылке
```bash
curl -i http://localhost:8080/<short_id>
```

### Удалить ссылку
```bash
curl -i -X DELETE http://localhost:8080/<short_id>
```

---

## Ошибки

### 400 Bad Request
Некорректный ввод (например, неподдерживаемая схема URL или неверный формат TTL).
```json
{
  "$schema": "/schemas/ErrorModel.json",
  "title": "Bad Request",
  "status": 400,
  "detail": "original_url must use http or https"
}
```

### 404 Not Found
Ссылка не найдена.
```json
{
  "$schema": "/schemas/ErrorModel.json",
  "title": "Not Found",
  "status": 404,
  "detail": "short_id not found"
}
```

### 409 Conflict
Коллизия `short_id`.
```json
{
  "$schema": "/schemas/ErrorModel.json",
  "title": "Conflict",
  "status": 409,
  "detail": "short_id conflict"
}
```

### 410 Gone
Ссылка больше недоступна (истёк срок жизни или превышен лимит переходов).
```json
{
  "$schema": "/schemas/ErrorModel.json",
  "title": "Gone",
  "status": 410,
  "detail": "expired"
}
```
