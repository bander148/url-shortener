package storage

import (
	"sync"
	"url-shortener/internal/models"
	"url-shortener/internal/util/shortcode"
)

type MemoryStorage struct {
	urls    map[string]models.UrlData
	mu      sync.RWMutex
	counter int64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls:    make(map[string]models.UrlData),
		counter: 1,
	}
}
func (m *MemoryStorage) Save(data models.UrlData) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	counter := m.counter
	m.counter++
	ShortKey := shortcode.ToBase62(counter)
	m.urls[ShortKey] = data
	return ShortKey, nil
}
func (m *MemoryStorage) Get(ShortKey string) (models.UrlData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.urls[ShortKey]
	if ok {
		return data, nil
	}
	return models.UrlData{}, ErrNotFound
}
func (m *MemoryStorage) Close() error {
	return nil
}
