package storage

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/util/shortcode"

	"github.com/redis/go-redis/v9"
)

// Структура для хранения подключения к Redis
// client — основной объект для взаимодействия с Redis.
// Под капотом управляет пулом TCP-соединений, автоматически их переоткрывает, балансирует нагрузку и т.д.
type RedisStorage struct {
	client *redis.Client
}

// Конфиг редис, возможно расширение для большего контроля
type Config struct {
	Addr     string `yaml:"addr"`     // Адрес, пример localhost:6379
	Password string `yaml:"password"` // Пароль
	// User        string        `yaml:"user"`  // User - имя пользователя
	DB          int           `yaml:"db"`           // идентификатор базы данных
	MaxRetries  int           `yaml:"max_retries"`  // максимальное количество попыток подключения
	DialTimeout time.Duration `yaml:"dial_timeout"` //таймаут для установления новых соединений
	Timeout     time.Duration `yaml:"timeout"`      // таймаут для записи и чтения.
}

// Создает новую RedisStorage
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// Save — сохраняет длинную ссылку и возвращает сгенерированный короткий ключ.
// Алгоритм:
// 1. Атомарно увеличиваем глобальный счётчик в Redis (ключ "shortener:counter")
// 2. Полученное число переводим в base62 → это и есть короткий ключ
// 3. Сохраняем пару shortKey → longUrl с TTL 30 дней
func (r *RedisStorage) Save(data models.UrlData) (string, error) {
	// Создание контекста
	ctx := context.Background()
	// Проверяем, есть ли уже такая длинная ссылка
	existingShortKey, err := r.client.HGet(ctx, "longurls", data.LongUrl).Result()
	if err == nil && existingShortKey != "" {
		// Уже существует — возвращаем существующий ключ
		return existingShortKey, nil
	}
	if err != redis.Nil && err != nil {
		return "", fmt.Errorf("hget longurls: %w", err)
	}
	// Генерируем уникальный числовой ID через атомарный инкремент в Redis
	// Ключ создаётся автоматически, если его ещё нет (INCR возвращает 1)
	// INCR — атомарная операция, гарантирует уникальность даже при высокой конкуренции
	counter, err := r.client.Incr(ctx, "shortener:counter").Result()
	if err != nil {
		return "", fmt.Errorf("incr counter: %w", err)
	}
	// Преобразуем числовой ID в короткую строку (base62)
	shortKey := shortcode.ToBase62(counter)

	// Сохраняем сопоставление shortKey → длинная ссылка, с TTL 30 дней
	// TTL = 30 суток — чтобы не засорять Redis ссылками, которые уже никто не использует
	err = r.client.Set(ctx, shortKey, data.LongUrl, 30*24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("set url: %w", err)
	}
	err = r.client.HSet(ctx, "longurls", data.LongUrl, shortKey).Err()
	if err != nil {
		return "", fmt.Errorf("hset longurls: %w", err)
	}
	return shortKey, nil
}

// Get — получает длинную ссылку по короткому ключу.
// Возвращает ErrNotFound, если ключа нет или он истёк.
func (r *RedisStorage) Get(shortKey string) (models.UrlData, error) {
	//Создание контекста
	ctx := context.Background()
	//Запрос к Redis по ключу shortKey
	val, err := r.client.Get(ctx, shortKey).Result()
	//Проверка есть ли нужное значение в Redis
	if err == redis.Nil {
		return models.UrlData{}, ErrNotFound // специальная ошибка для 404-сценария
	}
	//Проверка есть ли ошибки Redis
	if err != nil {
		return models.UrlData{}, fmt.Errorf("get key: %w", err)
	}
	//Возвращение значения согласно структуре
	return models.UrlData{LongUrl: val}, nil
}

// Закрытие подключения
func (r *RedisStorage) Close() error {
	return r.client.Close()
}

// NewClient — создаёт и проверяет соединение с Redis.
// Важно: клиент thread-safe и сам управляет пулом соединений — не нужно создавать новый на каждый запрос.
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})
	// Проверяем, что Redis действительно живой
	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}
	//Возвращаем указателль на redis.Client
	return db, nil
}
