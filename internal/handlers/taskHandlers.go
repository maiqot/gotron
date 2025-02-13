package handlers

import (
	"context"
	"firstProject/internal/tasksService" // Импортируем наш сервис
	"firstProject/internal/web/tasks"
	"fmt"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *tasksService.TaskService
}

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *tasksService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

// GetUsersIdTasks — получение задач по ID пользователя
func (h *Handler) GetUsersIdTasks(ctx context.Context, request tasks.GetUsersIdTasksRequestObject) (tasks.GetUsersIdTasksResponseObject, error) {
	// Получаем задачи пользователя по ID
	userTasks, err := h.Service.GetTasksByUserID(uint(request.Id))
	if err != nil {
		return nil, err
	}

	// Преобразуем задачи в формат API
	var response tasks.GetUsersIdTasks200JSONResponse
	for _, task := range userTasks {
		t := tasks.Task{
			Id:     &task.ID,
			Task:   &task.Task,
			IsDone: &task.IsDone,
			UserId: &task.UserID,
		}
		response = append(response, t)
	}

	// Возвращаем ответ в правильном формате
	return response, nil
}

func (h *Handler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	// Получение всех задач из сервиса
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	// Создаем переменную респон типа 200джейсонРеспонс
	// Которую мы потом передадим в качестве ответа
	response := tasks.GetTasks200JSONResponse{}

	// Заполняем слайс response всеми задачами из БД
	for _, task := range allTasks {
		task := tasks.Task{
			Id:     &task.ID,
			Task:   &task.Task,
			IsDone: &task.IsDone,
		}
		response = append(response, task)
	}

	// САМОЕ ПРЕКРАСНОЕ. Возвращаем просто респонс и nil!
	return response, nil
}

func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Отладочный вывод для проверки входных данных
	fmt.Printf("Received Task: %+v\n", request.Body)

	// Проверка, что тело запроса не пустое
	if request.Body == nil {
		return nil, echo.NewHTTPError(400, "Request body is required")
	}

	// Проверка, что поля Task и IsDone не пустые
	if request.Body.Task == nil || request.Body.IsDone == nil {
		return nil, echo.NewHTTPError(400, "Task and IsDone fields are required")
	}

	// Создаем задачу
	taskToCreate := tasksService.Task{
		Task:   *request.Body.Task,
		IsDone: *request.Body.IsDone,
	}

	// Сохраняем задачу через сервис
	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to create task")
	}

	// Формируем ответ
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Task,
		IsDone: &createdTask.IsDone,
	}

	// Возвращаем успешный ответ
	return response, nil
}

func (h *Handler) PatchTasksId(_ context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	id := uint(request.Id)
	updatedTask := tasksService.Task{
		Task:   *request.Body.Task,
		IsDone: *request.Body.IsDone,
	}

	// Вызываем метод сервиса
	task, err := h.Service.UpdateTaskByID(id, updatedTask)
	if err != nil {
		return nil, err
	}

	// Формируем ответ
	response := tasks.PatchTasksId200JSONResponse{
		Id:     &task.ID,
		Task:   &task.Task,
		IsDone: &task.IsDone,
	}
	return response, nil
}

func (h *Handler) DeleteTasksId(_ context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	id := uint(request.Id)

	// Вызываем метод удаления из сервиса
	err := h.Service.DeleteTaskByID(id)
	if err != nil {
		return nil, err
	}

	// Возвращаем статус 204
	return tasks.DeleteTasksId204Response{}, nil
}
