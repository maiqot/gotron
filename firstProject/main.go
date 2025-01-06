package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// task: глобальная переменная, которая будет хранить текущее значение задачи.
var task string

// RequestBody: структура, которая будет использоваться для декодирования JSON из тела POST-запроса.
type RequestBody struct {
	Task string `json:"task"`
}

// HelloHandler: обработчик для маршрута /api/hello. Он возвращает приветственное сообщение с текущим значением
// переменной task.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", task)
}

// TaskHandler: обработчик для маршрута /api/task. Он принимает JSON с полем task, обновляет глобальную переменную task
// и возвращает сообщение об успешном обновлении.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	var RequestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task = RequestBody.Task
	fmt.Fprintln(w, "Task updated successfully!")
}

// main: основная функция, которая инициализирует маршрутизатор, регистрирует обработчики для маршрутов
// и запускает HTTP-сервер на порту 8080.
func main() {
	router := mux.NewRouter()
	// api/hello — обрабатывает GET-запросы и возвращает приветственное сообщение с текущим значением переменной task.
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	// api/task — обрабатывает POST-запросы, принимает JSON с полем task и обновляет глобальную переменную task.
	router.HandleFunc("/api/task", TasksHandler).Methods("POST")
	http.ListenAndServe(":8080", router)
}
