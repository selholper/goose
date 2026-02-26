# REST API — трёхслойная архитектура + Goose миграции

Простой CRUD REST API на Go с трёхслойной архитектурой и автоматической накаткой миграций в PostgreSQL через [Goose](https://github.com/pressly/goose).

## Стек технологий

| Компонент       | Технология                        |
|-----------------|-----------------------------------|
| Язык            | Go 1.25                           |
| HTTP-роутер     | [chi](https://github.com/go-chi/chi) |
| База данных     | PostgreSQL 16                     |
| Драйвер БД      | [pgx/v5](https://github.com/jackc/pgx) |
| Миграции        | [goose/v3](https://github.com/pressly/goose) |
| Контейнеризация | Docker + Docker Compose           |

---

## Архитектура проекта

```
.
├── main.go                          # Точка входа, инициализация, graceful shutdown
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── migrations/
│   ├── 00001_create_users_table.sql # Миграция: создание таблицы users
│   └── 00002_add_users_index.sql    # Миграция: индекс на email
└── internal/
    ├── config/
    │   └── config.go                # Конфигурация из env-переменных
    ├── domain/
    │   └── user.go                  # Доменные модели (User, запросы)
    ├── repository/                  # Слой 1: работа с БД (pgx)
    │   └── user_repository.go
    ├── service/                     # Слой 2: бизнес-логика, валидация
    │   └── user_service.go
    └── handler/                     # Слой 3: HTTP-обработчики (chi)
        └── user_handler.go
```

### Три слоя

```
HTTP Request
    │
    ▼
┌─────────────┐
│   Handler   │  — декодирует JSON, вызывает Service, возвращает HTTP-ответ
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  — бизнес-логика, валидация, нормализация данных
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository  │  — SQL-запросы к PostgreSQL через pgx
└─────────────┘
```

---

## Быстрый старт (Docker Compose)

### 1. Клонировать репозиторий

```bash
git clone https://github.com/example/rest-api.git
cd rest-api
```

### 2. Запустить всё одной командой

```bash
docker compose up --build
```

При старте приложение **автоматически накатит миграции** через Goose и запустит сервер на порту `8080`.

Ожидаемый вывод:

```
restapi_app  | time=... level=INFO msg="connected to database" host=postgres db=restapi
restapi_app  | time=... level=INFO msg="OK   00001_create_users_table.sql"
restapi_app  | time=... level=INFO msg="OK   00002_add_users_index.sql"
restapi_app  | time=... level=INFO msg="migrations applied successfully"
restapi_app  | time=... level=INFO msg="server started" port=8080
```

---

## Запуск локально (без Docker)

### 1. Поднять PostgreSQL

```bash
docker run -d \
  --name restapi_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=restapi \
  -p 5432:5432 \
  postgres:16-alpine
```

### 2. Установить зависимости

```bash
go mod tidy
```

### 3. Запустить приложение

```bash
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=restapi \
HTTP_PORT=8080 \
go run main.go
```

Миграции накатятся автоматически при старте.

---

## Миграции Goose

Миграции хранятся в директории `migrations/` и встроены в бинарник через `//go:embed`.

### Формат файла миграции

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (...);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

### Ручное управление миграциями (CLI)

```bash
# Установить goose CLI
go install github.com/pressly/goose/v3/cmd/goose@latest

# Посмотреть статус миграций
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/restapi?sslmode=disable" status

# Применить все миграции
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/restapi?sslmode=disable" up

# Откатить последнюю миграцию
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/restapi?sslmode=disable" down

# Создать новую миграцию
goose -dir migrations create add_phone_to_users sql
```

---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

| Метод    | Путь          | Описание                  |
|----------|---------------|---------------------------|
| `GET`    | `/health`     | Проверка состояния сервиса |
| `POST`   | `/users`      | Создать пользователя       |
| `GET`    | `/users`      | Получить всех пользователей |
| `GET`    | `/users/{id}` | Получить пользователя по ID |
| `PUT`    | `/users/{id}` | Обновить пользователя      |
| `DELETE` | `/users/{id}` | Удалить пользователя       |

---

## Примеры запросов

### HTTP-клиент файл (`api/requests.http`)

В директории `api/` находится файл `requests.http` с готовыми примерами всех запросов.
Он поддерживается напрямую в:

- **JetBrains IDE** (IntelliJ IDEA, GoLand, WebStorm и др.) — встроенный HTTP Client
- **VS Code** — расширение [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)

Для запуска откройте `api/requests.http` в IDE и нажмите кнопку ▶ рядом с нужным запросом.

---

### Примеры через curl

### Проверить здоровье сервиса

```bash
curl http://localhost:8080/health
```

```json
{"status":"ok"}
```

### Создать пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Иван Иванов", "email": "ivan@example.com"}'
```

```json
{
  "id": 1,
  "name": "Иван Иванов",
  "email": "ivan@example.com",
  "created_at": "2026-02-26T08:00:00Z",
  "updated_at": "2026-02-26T08:00:00Z"
}
```

### Получить всех пользователей

```bash
curl http://localhost:8080/api/v1/users
```

```json
[
  {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@example.com",
    "created_at": "2026-02-26T08:00:00Z",
    "updated_at": "2026-02-26T08:00:00Z"
  }
]
```

### Получить пользователя по ID

```bash
curl http://localhost:8080/api/v1/users/1
```

### Обновить пользователя

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Пётр Петров", "email": "petr@example.com"}'
```

### Удалить пользователя

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

Успешный ответ: `204 No Content`

---

## Переменные окружения

| Переменная    | По умолчанию | Описание                    |
|---------------|--------------|-----------------------------|
| `HTTP_PORT`   | `8080`       | Порт HTTP-сервера            |
| `DB_HOST`     | `localhost`  | Хост PostgreSQL              |
| `DB_PORT`     | `5432`       | Порт PostgreSQL              |
| `DB_USER`     | `postgres`   | Пользователь БД              |
| `DB_PASSWORD` | `postgres`   | Пароль БД                    |
| `DB_NAME`     | `restapi`    | Имя базы данных              |
| `DB_SSLMODE`  | `disable`    | Режим SSL (`disable`/`require`) |

---

## HTTP коды ответов

| Код | Описание                              |
|-----|---------------------------------------|
| 200 | OK — запрос выполнен успешно          |
| 201 | Created — ресурс создан               |
| 204 | No Content — удаление успешно         |
| 400 | Bad Request — невалидный JSON         |
| 404 | Not Found — пользователь не найден    |
| 422 | Unprocessable Entity — ошибка валидации |
| 500 | Internal Server Error — ошибка сервера |
