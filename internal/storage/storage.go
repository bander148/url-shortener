package storage

import (
	"errors"
	"url-shortener/internal/models"
)

// Стандартная ошибка если url не найден
var ErrNotFound = errors.New("short key not found")

// Общий интерфейс для хранилища(БД)
type Storage interface {
	Save(data models.UrlData) (string, error)
	Get(shortKey string) (models.UrlData, error)
	Close() error
}
