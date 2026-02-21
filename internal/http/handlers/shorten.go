package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
)

// ShortenRequest — структура для входящего JSON от клиента
type ShortenRequest struct {
	LongURL string `json:"long_url"` // ← обязательно тег json:"long_url"
}

// ShortenResponse — структура для ответа клиенту
type ShortenResponse struct {
	ShortURL string `json:"short_url"` // ← тег json:"short_url"
}

// функция для получения длинной ссылки,создание короткой , ее сохранение, отправка пользователю
func ShortenHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	//Проверка является ли тип запроса POST
	if r.Method != http.MethodPost {
		//Выводим ошибку из-за неверного метода
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Объявляем переменную для входящего JSON от клиента
	var req ShortenRequest
	//Декодируем JSON в переменную req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//Если что то пошло не так выдаем ошибку
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//Проверяем нужно нам поле LongURL, если оно пустое выдаем ошибку
	if req.LongURL == "" {
		http.Error(w, "long_url is required", http.StatusBadRequest)
		return
	}
	//Проверка, является ли строка, которую прислал клиент в поле long_url, действительно корректным URL
	//(а не просто какой-то текст).
	parsedURL, err := url.ParseRequestURI(req.LongURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	//Создание данных согласно схеме UrlData
	data := models.UrlData{
		LongUrl: req.LongURL,
	}
	//Сохранение данных
	shortKey, err := store.Save(data)
	if err != nil {
		log.Printf("Save error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	//Создание короткой ссылки
	ShortUrl := fmt.Sprintf("http://localhost:8080/%s", shortKey)
	//Создание JSON для ответа
	response := ShortenResponse{ShortURL: ShortUrl}
	//Устанавливаем заголовки и говорми что хотим отправить JSON
	w.Header().Set("Content-Type", "application/json")
	//Устанавливает HTTP-статус-код ответа на 201 Created.
	w.WriteHeader(http.StatusCreated)
	//Берёт структуру response и превращает её в JSON-байты, которые клиент получит.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Encode error: %v", err)
	}
}
