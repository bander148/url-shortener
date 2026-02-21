# URL Shortener на Go

Простой сервис сокращения ссылок с поддержкой Redis.

## Запуск
go run cmd/app/main.go

## Примеры
POST http://localhost:8080/shorten
{"long_url": "https://google.com"}

GET http://localhost:8080/1 → редирект на google.com
