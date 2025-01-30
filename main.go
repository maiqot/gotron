package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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
		http.Error(w, "Ошибка при сохранении в БД", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task обновлен: %s ", reqBody.Message)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем ID задачи из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Декодируем JSON из тела запроса
	var reqBody struct {
		Task   string `json:"task,omitempty"`
		IsDone bool   `json:"is_done,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Находим задачу по ID
	var message Message
	result := DB.First(&message, id)
	if result.Error != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Обновляем поля задачи
	if reqBody.Task != "" {
		message.Task = reqBody.Task
	}
	if reqBody.IsDone != message.IsDone {
		message.IsDone = reqBody.IsDone
	}

	// Сохраняем изменения в БД
	DB.Save(&message)

	// Формируем ответ
	response := TaskResponse{
		ID:     message.ID,
		Task:   message.Task,
		IsDone: message.IsDone,
	}
	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задачи из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Находим задачу по ID
	var message Message
	result := DB.First(&message, id)
	if result.Error != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Удаляем задачу (мягкое удаление)
	DB.Delete(&message)

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Задача успешно удалена",
		"id":      id,
	})
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
	router.HandleFunc("/api/task/{id}", UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/api/delete/{id}", DeleteTaskHandler).Methods("DELETE")

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe("localhost:8080", router)
}
