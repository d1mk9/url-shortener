# URL Shortener

Минимальный сервис для сокращения URL и управления короткими ссылками.

## Возможности
- Создание короткой ссылки (`POST /shorten`)
- Переход по короткой ссылке (`GET /{short_id}`) — 307 Redirect
- Удаление ссылки (`DELETE /{short_id}`)
- Health‑check (`GET /healthz`)

---

## Установка и запуск

### Требования
- Go 1.22+
- PostgreSQL

### Настройки

Конфигурация задаётся через YAML‑файл и переменные окружения.

#### CONFIG_FILE (YAML)
Укажите путь к YAML‑файлу в переменной окружения `CONFIG_FILE`.

```yaml
# conf/config.yaml
postgres_host: "localhost"
postgres_port: "5432"
postgres_db: "urlshortener"
base_url: "http://urlshort.xy"
```

#### Важно про `base_url`
- Должен быть **абсолютным URL** (со схемой `http` или `https`).
- **Без завершающего слэша**. Пример: `http://urlshort.xy`.
- Используется для построения поля `ShortURL` в ответе (`<base_url>/<short_id>`).  
  Например, при `base_url: http://urlshort.xy` и `ShortID=abc123` сервис вернёт
  `ShortURL: http://urlshort.xy/abc123`.

#### Переменные окружения (секреты не коммитим)
```bash
export POSTGRES_USER=<username>
export POSTGRES_PASSWORD=<password>
export CONFIG_FILE=./conf/config.yaml
```

---

## Миграции

```bash
# применить
go run ./cmd/app migrate up

# откатить одну
go run ./cmd/app migrate down
```

---

## Запуск сервера

```bash
POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres CONFIG_FILE=./conf/config.yaml   go run ./cmd/app serve
```

Сервер поднимется на `http://localhost:8080`.

---

## Makefile (опционально)

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

---

## REST API

### Создать короткую ссылку
- `Expiry` — строка формата `time.Duration` (например, `24h`, `15m`). Опционально.  
- `MaxVisits` — максимум переходов. `0` означает «без ограничений». Опционально.

```bash
curl -i -X POST http://localhost:8080/shorten   -H "Content-Type: application/json"   -d '{"OriginalURL":"https://example.com","MaxVisits":3,"Expiry":"24h"}'
```

Пример ответа:
```json
{
  "$schema": "http://localhost:8080/schemas/ShortenResponse.json",
  "ShortURL": "http://urlshort.xy/abc123xy",
  "ShortID": "abc123xy",
  "ExpiresAt": "2025-09-15T12:00:00Z",
  "MaxVisits": 3
}
```

### Перейти по короткой ссылке
```bash
curl -i http://localhost:8080/<short_id>
```
Ответ: `307 Temporary Redirect` + заголовок `Location: <original_url>` и `Cache-Control: no-store`.

### Удалить ссылку
```bash
curl -i -X DELETE http://localhost:8080/<short_id>
```

---

## Ошибки

### 400 Bad Request
Некорректный ввод (неподдерживаемая схема URL, неверный формат `Expiry`).
```json
{
  "$schema": "/schemas/ErrorModel.json",
  "title": "Bad Request",
  "status": 400,
  "detail": "invalid expiry duration"
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
Коллизия при генерации `short_id`.
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

---
