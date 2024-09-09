package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	storage Storage
	dbFlag  string
)

func main() {
	// чтение флага
	flag.StringVar(&dbFlag, "d", "memory", "Select storage: memory or postgres")
	flag.Parse()

	if dbFlag == "postgres" {
		storage = NewPostgresStorage()
	} else {
		storage = NewInMemoryStorage()
	}

	// вызов функций
	http.HandleFunc("/", handleShortenURL)
	http.HandleFunc("/get/", handleGetURL)

	// Запуск сервера
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
