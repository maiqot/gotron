package main

import (
	"encoding/json"
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
	var messages []Message
	// Извлекаем все записи из базы данных
	result := DB.Find(&messages)
	if result.Error != nil {
		http.Error(w, "Не удалось получить сообщения", http.StatusInternalServerError)
		return
	}
	// Формируем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var RequestBody requestBody
	// Декодируем JSON из тела запроса в requestBody
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		http.Error(w, "Неправильный формат JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что task не пустое
	if RequestBody.Task == "" {
		http.Error(w, "Поле task не может быть пустым", http.StatusBadRequest)
		return
	}

	// Создаем объект Message для сохранения в БД
	message := Message{
		Task:   RequestBody.Task, // Передаем поле task из JSON
		IsDone: false,            // По умолчанию устанавливаем IsDone в false
	}

	// Сохраняем объект в БД
	result := DB.Create(&message)
	if result.Error != nil {
		http.Error(w, "Не удалось сохранить сообщение в базе данных", http.StatusInternalServerError)
		return
	}
	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
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
