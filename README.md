# Url_sh

пример использование
- запуск
- go run main.go storage.go utils.go handlers.go -d (memory(для сохранения в оперативную память)/postgres)
- после обрабатываем запросы мне удобнее через curl
- curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -d "{\"url\":\"http://example.com\"}"
- curl -X GET http://localhost:8080/get/shortURL
