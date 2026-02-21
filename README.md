# URL Shortener на Go

Простой и надёжный сервис сокращения ссылок, написанный на Go.  
Поддерживает хранение в памяти и Redis, автоматическую генерацию коротких ключей (counter + base62), проверку дубликатов длинных ссылок, валидацию URL и graceful shutdown.

## Возможности

- Создание короткой ссылки (POST /shorten)
- Редирект по короткой ссылке (GET /{shortKey})
- Хранение в памяти (для простоты) или Redis (production-ready)
- Проверка дубликатов длинных ссылок (не создаёт лишние ключи)
- Генерация коротких ключей без коллизий (base62 от счётчика)
- Graceful shutdown (Ctrl+C → корректное завершение)
- Валидация URL и обработка ошибок

## Технологии

- Go 1.22+
- Redis (go-redis/v9)
- net/http (встроенный роутер)
- sync.RWMutex (для in-memory)
- base62 (собственная реализация)

## Быстрый запуск

### 1. Клонируй репозиторий

```bash
git clone https://github.com/bander148/url-shortener.git
cd url-shortener
```
### 2. Запуск в памяти (без Redis) — самый простой вариант
```bash
Bashgo run cmd/app/main.go
```
Сервер запустится на:
http://localhost:8080

### 3. Запуск с Redis (рекомендуется для production-подобного поведения)

Вариант A — через Docker (самый быстрый)
```bash
Bashdocker run -d -p 6379:6379 --name redis-local redis:latest
```
Вариант B — без Docker
Установи Redis для Windows (Memurai или нативная сборка):
→ https://www.memurai.com/get-memurai (Developer Edition — бесплатно)

Или используй WSL / Linux.

После запуска Redis выполни:
```bash
Bashgo run cmd/app/main.go
```
Сервер снова на http://localhost:8080

### Примеры запросов (Postman / curl)
Создать короткую ссылку
```text
Bashcurl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://google.com"}'
Ответ (201 Created):
JSON{
  "short_url": "http://localhost:8080/1"
}
```

Редирект по короткой ссылке

Открой в браузере:
http://localhost:8080/1
→ мгновенно перенаправит на https://google.com

### Проверка ошибок

Несуществующий ключ

http://localhost:8080/999 → 404 Not Found
```text
Невалидный URLBashcurl
 -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "google.com"}'→ 400 Bad Request ("Invalid URL")

Пустой long_urlBashcurl
 -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": ""}'→ 400 Bad Request ("long_url is required")

```

### Структура проекта 
```text
url-shortener/
├── cmd/
│   └── app/
│       └── main.go           # точка входа + graceful shutdown
├── internal/
│   ├── handlers/             # HTTP-обработчики
│   │   ├── shorten.go        # POST /shorten — создание короткой ссылки
│   │   └── redirect.go       # GET /{key} — редирект или 404
│   ├── storage/              # хранилища (интерфейс + реализации)
│   │   ├── storage.go        # интерфейс Storage
│   │   ├── memory.go         # реализация в памяти (с мьютексом)
│   │   └── redis.go          # реализация на Redis (go-redis/v9)
│   └── models/               # общие структуры данных
│       └── url_data.go
├── go.mod                    # зависимости
└── README.md
```
### Следующие шаги / улучшения
- Тесты на handlers (httptest + table-driven)
- Счётчик кликов по каждой ссылке
- TTL для ссылок (автоудаление через 30 дней)
- Деплой на Render / Fly.io / Railway
- Middleware (rate limiting, logging, CORS)
- Конкурентность (worker pool для обработки кликов)
