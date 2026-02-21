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
2. Запусти в памяти (без Redis)
Bashgo run cmd/app/main.go
Сервер запустится на http://localhost:8080
3. Запусти с Redis (рекомендуется)
Запусти Redis (самый простой способ — Docker):
Bashdocker run -d -p 6379:6379 --name redis-local redis:latest
Или используй Memurai/нативный Redis для Windows.
Затем запусти проект:
Bashgo run cmd/app/main.go
Примеры запросов (Postman / curl)
Создать короткую ссылку
Bashcurl -X POST -H "Content-Type: application/json" \
-d '{"long_url":"https://google.com"}' \
http://localhost:8080/shorten
Ответ (201 Created):
JSON{"short_url":"http://localhost:8080/1"}
Редирект
Открой в браузере:
http://localhost:8080/1
→ мгновенно перенаправит на https://google.com
Несуществующий ключ
http://localhost:8080/999 → 404 Not Found
Структура проекта
texturl-shortener/
├── cmd/
│   └── app/
│       └── main.go           # точка входа + graceful shutdown
├── internal/
│   ├── handlers/             # HTTP-обработчики
│   │   ├── shorten.go
│   │   └── redirect.go
│   ├── storage/              # хранилища (интерфейс + реализации)
│   │   ├── storage.go        # интерфейс Storage
│   │   ├── memory.go
│   │   └── redis.go
│   └── models/               # структуры данных
│       └── url_data.go
├── go.mod
└── README.md
Следующие шаги / улучшения

Тесты на handlers (httptest + table-driven)
Счётчик кликов по каждой ссылке
TTL для ссылок (автоудаление через 30 дней)
Деплой на Render / Fly.io / Railway
Middleware (rate limiting, logging, CORS)
Конкурентность (worker pool для обработки кликов)
