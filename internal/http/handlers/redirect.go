package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"url-shortener/internal/storage"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	//Проверить Get ли запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Спарсить из request ShortKey
	shortKey := strings.Trim(r.URL.Path, "/")
	//Вызвать метод GET
	data, err := store.Get(shortKey)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Short link not found", http.StatusNotFound)
		} else {
			log.Printf("Get error : %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	//Получить нужную longurl
	longUrl := data.LongUrl
	//Перенаправить пользователя по ссылке(не забыть статус код перенаправления)
	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}
