# 🚀 GitHub Release Notifier

**GitHub Release Notifier** — це сучасний мікросервіс на Go, розроблений для моніторингу оновлень у GitHub репозиторіях та автоматичного сповіщення користувачів. Проект побудований за принципами **Clean Architecture** та **SOLID**.

---

## ✨ Ключові особливості

* **Smart Caching**: Реалізовано паттерн **Decorator** для кешування запитів до GitHub API через **Redis**, що значно економить ліміти токена.
* **Full Observability**: Повний моніторинг стану системи через **Prometheus** (метрики Gin, Go Runtime та бізнес-показники).
* **Rate Limit Protection**: Адаптивне керування лімітами GitHub API — сервіс автоматично "засинає" до моменту скидання ліміту.
* **Interactive UI**: Вбудована сучасна HTML-сторінка (Tailwind CSS) для швидкої підписки.
* **API Documentation**: Повна підтримка **Swagger (OpenAPI 3.0)** для зручного тестування ендпоінтів.

---

## 🛠 Технологічний стек

* **Мова**: Go 1.26 (Gin Gonic framework)
* **База даних**: PostgreSQL 17 (pgx pool)
* **Кешування**: Redis 7
* **Моніторинг**: Prometheus
* **Воркери**: `gocron/v2` для планування сканування.
* **Тестування**: `testify` (Table-driven tests, Mocking).

---

## 📂 Структура проекту

Проект організований згідно з принципами чистої архітектури:

* `cmd/` — точка входу, ініціалізація та Dependency Injection.
* `internal/domain/` — бізнес-сутності та інтерфейси.
* `internal/usecase/` — ядро бізнес-логіки (покрито тестами на **82.6%**).
* `internal/infrastructure/` — зовнішні інтеграції (DB, GitHub API з Redis Decorator, Email).
* `internal/delivery/` — HTTP хендлери та Middleware (Auth, Prometheus).
* `web/` — фронтенд частина (HTML/JS).

---

## 🚀 Швидкий запуск (Docker)

### 1. Налаштування середовища
Створіть файл `.env` у корені проекту (використовуйте `.env.example` як шаблон):

```env
ENV=development
HOST=0.0.0.0
PORT=8080
SCANNER_INTERVAL=5m
API_KEY=grn-subscription-service-secret-key

#SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=
SMTP_FROM=noreply@grn-subscription-service.com

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASS=
REDIS_DB=0
REDIS_CACHE_TTL=10m

#Database
POSTGRES_DSN=postgresql://dev:devpassv2@db:5432/grn-subscription-service-db?sslmode=disable

# Database Configuration
DB_USER=dev
DB_NAME=grn-subscription-service-db
DB_PASS_CONTAINER_PATH=/run/secrets/db_pass_secret

GITHUB_TOKEN=

```

---

## 2. Запуск всієї інфраструктури

```bash
docker-compose up --build
```
---

## 🔍 Моніторинг та Документація
### Після запуску сервіси доступні за наступними адресами:

* 🏠 Головна сторінка: http://localhost:8080/

* 📜 Swagger UI: http://localhost:8080/swagger/index.html

* 📈 Prometheus UI: http://localhost:9090

* 📊 Metrics Endpoint: http://localhost:8080/metrics

## 🧪 Тестування
### Проект використовує ідіоматичні для Go Table-driven tests.

## Запуск тестів з перевіркою покриття:

```bash
go test ./internal/usecase/... -coverprofile=cover.out
go tool cover -func=cover.out
```
