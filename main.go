package main

import (
	"fmt"
	"log"
	"net/http"
)

// Обработчик запроса по пути /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	// Регистрируем обработчик для пути /health
	http.HandleFunc("/health", healthHandler)
	port := "8080"
	log.Printf("Запуск сервиса на порту %s...\n", port)
	// Запускаем сервер
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
