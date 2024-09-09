package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

// Обработчик запросов для сокращения
func handleShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Декодирование JSON-запроса
	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация URL(проверка на пустоту)
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Генерация сокращенного URL
	shortURL := generateShortURL()

	// Сохранение оригинального URL в хранилище
	err = storage.SaveURL(shortURL, req.URL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	// Ответ с сокращённым URL
	responseURL := "http://localhost:8080/get/" + shortURL
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseURL))
}

// Обработчик GET-запросов
func handleGetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение короткого URL из пути запроса
	shortURL := strings.TrimPrefix(r.URL.Path, "/get/")

	// Проверка наличия короткого URL
	if shortURL == "" {
		http.Error(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	// Получение оригинального URL из хранилища
	originalURL, err := storage.GetURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Перенаправление на оригинальный URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}
