package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/http/handlers"
	"url-shortener/internal/storage"
)

func main() {
	//Создание конфига (жедательно вынести в отдельный файл и получать данные из .env)
	cfg := storage.Config{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		DialTimeout: 5 * time.Second,
		Timeout:     3 * time.Second,
		MaxRetries:  3,
	}
	// 1. Создаём хранилище
	redisClient, err := storage.NewClient(context.Background(), cfg)

	if err != nil {
		log.Fatalf("redis connect: %v", err)
	}
	defer redisClient.Close()
	st := storage.NewRedisStorage(redisClient)

	// 2. Создаём мультиплексор и привязываем обработчики
	mux := http.NewServeMux()
	// Хандлер для пост запросов с ссылками
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenHandler(w, r, st)
	})
	// Хандлер для перенаправление по ShortKey из URL
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, st)
	})
	// 3. Создаём сервер (не запускаем сразу)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// 4. Запускаем сервер в отдельной горутине
	go func() {
		log.Println("Сервер запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// 5. Ожидаем сигнал завершения (Ctrl+C или SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // здесь main блокируется, пока не придёт сигнал
	log.Println("Получен сигнал завершения. Запускаем graceful shutdown...")

	// 6. Даём серверу 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка graceful shutdown: %v", err)
	}

	// 7. Закрываем хранилище
	if err := st.Close(); err != nil {
		log.Printf("Ошибка закрытия хранилища: %v", err)
	}
	log.Println("Сервер остановлен корректно")

}
