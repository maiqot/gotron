package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	var messages []Message

	// Получаем все записи из БД
	result := DB.Find(&messages)
	if result.Error != nil {
		http.Error(w, "Ошибка при получении данных из БД", http.StatusInternalServerError)
		return
	}

	// Возвращаем данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody requestBody

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Создаем новую запись в БД
	message := Message{Task: reqBody.Message, IsDone: false}
	result := DB.Create(&message)
	if result.Error != nil {
		http.Error(w, "Ошибка при сохранение в БД", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task обновлен: %s ", reqBody.Message)
}

func main() {

	// Инициализация базы данных
	InitDB()

	// Автомиграция (создание таблицы, если она не существует)
	DB.AutoMigrate(&Message{})

	// Настройка маршрутов
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/api/post", PostHandler).Methods("POST")

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe("localhost:8080", router)
}
