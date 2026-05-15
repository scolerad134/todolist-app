# TodoApp

Веб-приложение для управления задачами (to-do list) с REST API, веб-интерфейсом и аналитикой. Бэкенд написан на Go, данные хранятся в PostgreSQL. Одностраничный UI встроен в сервис и отдаётся с корневого маршрута `/`.

## Демо

Сервис развёрнут и доступен в интернете:

**[http://80.249.145.250:5050/](http://80.249.145.250:5050/)** — веб-интерфейс TodoApp

Документация API (Swagger UI): [http://80.249.145.250:5050/swagger/](http://80.249.145.250:5050/swagger/)

---

## Возможности

- **Пользователи** — создание, просмотр, частичное обновление (PATCH) и удаление; опциональный номер телефона в международном формате (`+7...`).
- **Задачи** — CRUD, привязка к автору (`author_user_id`), статусы «выполнена / не выполнена», даты создания и завершения.
- **Статистика** — количество созданных и выполненных задач, доля выполненных, среднее время выполнения; фильтры по пользователю и диапазону дат.
- **Веб-интерфейс** — адаптивная SPA-страница (`public/index.html`) с поддержкой светлой/тёмной темы.
- **REST API v1** — версионирование через префикс `/api/v1`.
- **Оптимистичная блокировка** — поле `version` у пользователей и задач; при конкурентном изменении возвращается `409 Conflict`.
- **Трёхсостоятельный PATCH** — поле можно не передавать (без изменений), передать значение или явно `null` (очистить, где допустимо).
- **Наблюдаемость** — структурированные логи (zap), `X-Request-ID`, трассировка запросов, обработка panic в middleware.

---

## Стек технологий

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.26 |
| HTTP | `net/http`, `ServeMux` (Go 1.22+ routing) |
| БД | PostgreSQL 17 |
| Драйвер БД | [pgx/v5](https://github.com/jackc/pgx) |
| Миграции | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Конфигурация | [envconfig](https://github.com/kelseyhightower/envconfig) |
| Валидация | [go-playground/validator](https://github.com/go-playground/validator) |
| Логирование | [uber-go/zap](https://github.com/uber-go/zap) |
| Документация API | [swaggo/swag](https://github.com/swaggo/swag) |
| Контейнеризация | Docker, Docker Compose |

---

## Архитектура

Проект организован по **feature-based** структуре с разделением слоёв:

```
transport (HTTP) → service → repository (postgres)
```

Общие сущности и инфраструктура вынесены в `internal/core`:

- `domain` — модели и бизнес-правила валидации
- `config`, `logger`, `errors`
- `transport/http` — сервер, middleware, коды ответов, утилиты
- `repository/postgres/pool` — пул соединений с БД

Каждая фича (`users`, `tasks`, `statistics`, `web`) содержит свой набор `service`, `repository`, `transport/http`.

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│  Browser    │────▶│  HTTP Server │────▶│  Feature        │
│  (index.html)│     │  + Middleware│     │  Services       │
└─────────────┘     └──────┬───────┘     └────────┬────────┘
                           │                      │
                           │              ┌───────▼────────┐
                           │              │  PostgreSQL    │
                           │              │  (schema       │
                           └──────────────│   todoapp)     │
                                          └────────────────┘
```

---

## Структура репозитория

```
todolist-app/
├── cmd/todoapp/           # Точка входа и Dockerfile приложения
├── internal/
│   ├── core/              # Общая инфраструктура
│   └── features/
│       ├── users/         # Пользователи
│       ├── tasks/         # Задачи
│       ├── statistics/    # Статистика
│       └── web/           # Отдача главной страницы
├── public/
│   └── index.html         # Веб-интерфейс
├── migrations/            # SQL-миграции (golang-migrate)
├── docs/                  # Сгенерированная Swagger-документация
├── docker-compose.yaml
├── Makefile
├── .env.examle            # Пример переменных окружения
└── go.mod
```

---

## Требования

- [Go](https://go.dev/dl/) 1.26+
- [Docker](https://www.docker.com/) и Docker Compose
- [Make](https://www.gnu.org/software/make/) (опционально, для команд из Makefile)

---

## Быстрый старт

### 1. Клонирование и настройка окружения

```bash
git clone https://github.com/scolerad134/todolist-app.git
cd todolist-app
cp .env.examle .env
```

Отредактируйте `.env`: задайте учётные данные PostgreSQL (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`).

### 2. Запуск PostgreSQL

```bash
make env-up
make migrate-up
```

Для локальной разработки приложения (без Docker-контейнера todoapp) пробросьте порт БД на хост:

```bash
make env-port-forward
```

PostgreSQL будет доступен на `localhost:5433`.

### 3. Запуск приложения

**Локально (go run):**

```bash
make todoapp-run
```

**В Docker (production-like):**

```bash
make todoapp-deploy
```

Приложение слушает порт **5050** (по умолчанию `HTTP_ADDR=:5050`).

| Ресурс | URL (локально) |
|--------|----------------|
| Веб-интерфейс | http://localhost:5050/ |
| API | http://localhost:5050/api/v1/ |
| Swagger UI | http://localhost:5050/swagger/ |

### 4. Остановка

```bash
# Остановить только приложение в Docker
make todoapp-undelpoy

# Остановить PostgreSQL
make env-down

# Закрыть проброс порта
make env-port-close
```

---

## Переменные окружения

Файл-пример: `.env.examle`. Makefile подключает `.env` через `include .env`.

| Переменная | Описание | Пример |
|------------|----------|--------|
| `HTTP_ADDR` | Адрес и порт HTTP-сервера | `:5050` |
| `HTTP_SHUTDOWN_TIMEOUT` | Таймаут graceful shutdown | `30s` |
| `HTTP_ALLOWED_ORIGINS` | Разрешённые Origin для CORS (через запятую) | `http://localhost:5050,null` |
| `POSTGRES_HOST` | Хост PostgreSQL | `localhost` / `todoapp-postgres` |
| `POSTGRES_PORT` | Порт PostgreSQL | `5432` / `5433` |
| `POSTGRES_USER` | Пользователь БД | — |
| `POSTGRES_PASSWORD` | Пароль БД | — |
| `POSTGRES_DB` | Имя базы данных | — |
| `POSTGRES_TIMEOUT` | Таймаут операций с БД | `10s` |
| `LOGGER_LEVEL` | Уровень логирования | `DEBUG` |
| `LOGGER_FOLDER` | Каталог для файлов логов | `./out/logs` |
| `TIME_ZONE` | Часовой пояс приложения | `UTC` |
| `PROJECT_ROOT` | Корень проекта (для пути к `public/`) | задаётся в Makefile |

Префиксы для `envconfig`: `HTTP_*`, `POSTGRES_*`, `LOGGER_*`.

---

## API

Базовый путь: **`/api/v1`**

### Пользователи

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/users` | Создать пользователя |
| `GET` | `/users` | Список пользователей (`limit`, `offset`) |
| `GET` | `/users/{id}` | Пользователь по ID |
| `PATCH` | `/users/{id}` | Частичное обновление |
| `DELETE` | `/users/{id}` | Удалить пользователя |

**Создание пользователя:**

```json
POST /api/v1/users
{
  "full_name": "Иван Иванов",
  "phone_number": "+79998887766"
}
```

- `full_name` — обязательно, 3–100 символов (по рунам).
- `phone_number` — опционально, 10–15 символов, формат `+[цифры]`.

### Задачи

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/tasks` | Создать задачу |
| `GET` | `/tasks` | Список задач (`user_id`, `limit`, `offset`) |
| `GET` | `/tasks/{id}` | Задача по ID |
| `PATCH` | `/tasks/{id}` | Частичное обновление |
| `DELETE` | `/tasks/{id}` | Удалить задачу |

**Создание задачи:**

```json
POST /api/v1/tasks
{
  "title": "Тренировка",
  "description": "Начало в 19:30",
  "author_user_id": 1
}
```

- `title` — обязательно, 1–100 символов.
- `description` — опционально, 1–1000 символов.
- `author_user_id` — ID существующего пользователя.

**PATCH (трёхсостоятельная логика):**

1. Поле **не передано** — значение в БД не меняется.
2. Поле передано со **значением** — устанавливается новое значение.
3. Поле передано как **`null`** — поле очищается (для `description`, `phone_number`). Для `full_name`, `title`, `completed` установка в `null` запрещена.

При одновременном изменении одной записи второй запрос может получить **`409 Conflict`** (оптимистичная блокировка по `version`).

### Статистика

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/statistics` | Агрегированная статистика по задачам |

Query-параметры (все опциональны):

- `user_id` — фильтр по автору задач.
- `from` — начало периода (включительно), формат `YYYY-MM-DD`.
- `to` — конец периода (не включительно), формат `YYYY-MM-DD`.

**Пример ответа:**

```json
{
  "tasks_created": 50,
  "tasks_completed": 10,
  "tasks_completed_rate": 20,
  "tasks_average_completion_time": "1m30s"
}
```

### Формат ошибок

```json
{
  "error": "полный текст ошибки",
  "message": "краткое сообщение для клиента"
}
```

Коды: `400` (невалидные данные), `404` (не найдено), `409` (конфликт версий), `500` (внутренняя ошибка).

---

## Swagger

Интерактивная документация доступна по адресу `/swagger/`.

Перегенерация после изменения аннотаций в коде:

```bash
make swagger-gen
```

Исходные аннотации — в `cmd/todoapp/main.go` и хендлерах `internal/features/*/transport/http/`.

---

## База данных

Схема: **`todoapp`**

### Таблица `users`

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `SERIAL` | Первичный ключ |
| `version` | `BIGINT` | Версия для optimistic locking |
| `full_name` | `VARCHAR(100)` | Имя, 3–100 символов |
| `phone_number` | `VARCHAR(15)` | Опционально, `+[0-9]`, 10–15 символов |

### Таблица `tasks`

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `SERIAL` | Первичный ключ |
| `version` | `BIGINT` | Версия для optimistic locking |
| `title` | `VARCHAR(100)` | Заголовок |
| `description` | `VARCHAR(100)` | Описание (опционально) |
| `completed` | `BOOLEAN` | Статус выполнения |
| `created_at` | `TIMESTAMPTZ` | Дата создания |
| `completed_at` | `TIMESTAMPTZ` | Дата завершения (если выполнена) |
| `author_user_id` | `INT` | FK → `users(id)` |

Миграции:

```bash
make migrate-up      # применить
make migrate-down    # откатить последнюю
make migrate-create seq=<имя>   # создать новую миграцию
```

---

## Makefile — справочник команд

| Команда | Назначение |
|---------|------------|
| `make env-up` | Запустить PostgreSQL |
| `make env-down` | Остановить PostgreSQL |
| `make env-cleanup` | Удалить volume с данными БД (с подтверждением) |
| `make env-port-forward` | Проброс PostgreSQL на `127.0.0.1:5433` |
| `make env-port-close` | Остановить проброс порта |
| `make migrate-up` / `migrate-down` | Миграции |
| `make migrate-create seq=<name>` | Новая миграция |
| `make todoapp-run` | Запуск через `go run` |
| `make todoapp-deploy` | Сборка и запуск в Docker |
| `make todoapp-undelpoy` | Остановить контейнер приложения |
| `make swagger-gen` | Генерация Swagger |
| `make logs_cleanup` | Очистить каталог логов |
| `make ps` | Статус контейнеров Compose |

---

## Docker Compose

Сервисы в `docker-compose.yaml`:

| Сервис | Назначение |
|--------|------------|
| `todoapp` | Приложение (порт `5050`) |
| `todoapp-postgres` | PostgreSQL 17 |
| `todoapp-postgres-migrate` | Одноразовый запуск миграций |
| `port-forwader` | Проброс БД на хост (`5433`) |
| `swagger` | Контейнер для `swag` CLI |

Сборка приложения — multi-stage Dockerfile в `cmd/todoapp/Dockerfile` (образ `golang:1.26` → `alpine:3.23`).

---

## Разработка

### Запуск без Docker (только БД в Docker)

```bash
make env-up
make env-port-forward
make migrate-up
make todoapp-run
```

### Логи

Логи пишутся в каталог из `LOGGER_FOLDER` (по умолчанию `out/logs`). Каталог в `.gitignore`.

### Middleware

- **CORS** — по списку `HTTP_ALLOWED_ORIGINS`
- **RequestID** — заголовок `X-Request-ID`
- **Logger** — контекстный логгер с `request_id` и URL
- **Trace** — логирование входящих запросов и латентности
- **Panic** — перехват panic → `500`

### Полезные curl-примеры

```bash
# Создать пользователя
curl -s -X POST http://localhost:5050/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Тест Тестов"}'

# Создать задачу
curl -s -X POST http://localhost:5050/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Первая задача","author_user_id":1}'

# Статистика
curl -s 'http://localhost:5050/api/v1/statistics?user_id=1'
```

---

## Лицензия

Уточните лицензию в репозитории при публикации. Модуль Go: `github.com/scolerad134/todolist-app`.
