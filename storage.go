package main

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

// Интерфейс для хранения URL
type Storage interface {
	SaveURL(shortURL, originalURL string) error
	GetURL(shortURL string) (string, error)
}

// ----------- Хранилище в памяти -----------

type InMemoryStorage struct {
	mu      sync.RWMutex
	storage map[string]string
}

// Конструктор нового хранилища в памяти
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[string]string),
	}
}

// Сохранение сокращенного URL в памяти
func (s *InMemoryStorage) SaveURL(shortURL, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[shortURL] = originalURL
	return nil
}

// Получение оригинального URL по сокращенному из памяти
func (s *InMemoryStorage) GetURL(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if originalURL, exists := s.storage[shortURL]; exists {
		return originalURL, nil
	}
	return "", errors.New("URL not found")
}

// ----------- Хранилище на базе bd -----------

type PostgresStorage struct {
	db *sql.DB
}

// Конструктор для создания нового хранилища PostgreSQL
func NewPostgresStorage() *PostgresStorage {
	connStr := "user=postgres password=IlS-yn100% dbname=url_shortener sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к PostgreSQL:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось выполнить проверку соединения с PostgreSQL:", err)
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            short_url VARCHAR(255) PRIMARY KEY,
            original_url TEXT NOT NULL
        )
    `)
	if err != nil {
		log.Fatal("Не удалось создать таблицу:", err)
	}
	return &PostgresStorage{db: db}
}

// Сохранение сокращенного URL в PostgreSQL
func (s *PostgresStorage) SaveURL(shortURL, originalURL string) error {
	_, err := s.db.Exec(
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING",
		shortURL, originalURL,
	)
	return err
}

// Получение оригинального URL по сокращенному из PostgreSQL
func (s *PostgresStorage) GetURL(shortURL string) (string, error) {
	var originalURL string

	// Выполняем запрос и проверяем на ошибки
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("URL not found")
		}
		return "", errors.New("query error: " + err.Error())

	}

	return originalURL, nil
}
