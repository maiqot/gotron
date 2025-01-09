package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// глобальная переменная для хранения значения task
var task string

// структура для декодирования JSON из тела POST–запроса
type requestBody struct {
	Task string `json:"task"`
}

// структура для формирования JSON-ответа
type responseBody struct {
	Message string `json:"message"`
	Task    string `json:"task"`
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", task)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var RequestBody requestBody
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Обновление глобальной переменной task
	task = RequestBody.Task

	// Возвращаем простой текстовый ответ
	fmt.Fprintf(w, "Задача обновлена: %s", task)
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	http.ListenAndServe(":8080", router)
}
