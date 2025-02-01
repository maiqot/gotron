package handlers

import (
	"encoding/json"
	"firstProject/internal/taskService" // Импортируем наш сервис
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *taskService.TaskService
}

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task taskService.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdTask, err := h.Service.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) PatchTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // Получаем параметры из URL
	id, err := strconv.Atoi(vars["id"]) // Преобразуем id в число
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Task   string `json:"task,omitempty"`
		IsDone bool   `json:"is_done,omitempty"`
	}

	// Декодируем JSON из тела запроса
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для обновления задачи
	updatedTask, err := h.Service.UpdateTaskByID(uint(id), taskService.Task{
		Task:   reqBody.Task,
		IsDone: reqBody.IsDone,
	})

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // Получаем параметры из URL
	id, err := strconv.Atoi(vars["id"]) // Преобразуем id в число
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для удаления задачи
	err = h.Service.DeleteTaskByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ об успешном удалении
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Задача успешно удалена",
		"id":      vars["id"],
	})
}
