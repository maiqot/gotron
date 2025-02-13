package handlers

import (
	"context"
	"errors"
	"firstProject/internal/tasksService"
	"firstProject/internal/web/tasks"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	Service *tasksService.TaskService
}

// NewHandler создает новый Handler
func NewHandler(service *tasksService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

// GetUsersIdTasks — получение задач по ID пользователя
func (h *Handler) GetUsersIdTasks(ctx context.Context, request tasks.GetUsersIdTasksRequestObject) (tasks.GetUsersIdTasksResponseObject, error) {
	userTasks, err := h.Service.GetTasksByUserID(uint(request.Id))
	if err != nil {
		return nil, err
	}

	var response tasks.GetUsersIdTasks200JSONResponse
	for _, task := range userTasks {
		id := uint(task.ID)         // Преобразуем int64 в uint
		userId := uint(task.UserID) // Преобразуем int64 в uint
		t := tasks.Task{
			Id:     &id,
			Task:   &task.Task,
			IsDone: &task.IsDone,
			UserId: &userId,
		}
		response = append(response, t)
	}

	return response, nil
}

// GetTasks — получение всех задач
func (h *Handler) GetTasks(ctx context.Context, request tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var response tasks.GetTasks200JSONResponse
	for _, task := range allTasks {
		id := uint(task.ID)         // Преобразуем int64 в uint
		userId := uint(task.UserID) // Преобразуем int64 в uint
		t := tasks.Task{
			Id:     &id,
			Task:   &task.Task,
			IsDone: &task.IsDone,
			UserId: &userId,
		}
		response = append(response, t)
	}

	return response, nil
}

// PostTasks — создание новой задачи
func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(400, "Request body is required")
	}

	if request.Body.Task == "" || request.Body.IsDone == nil || request.Body.UserId == 0 {
		return nil, echo.NewHTTPError(400, "Task, IsDone, and UserId fields are required")
	}

	taskToCreate := tasksService.Task{
		Task:   request.Body.Task,
		IsDone: *request.Body.IsDone,
		UserID: uint(request.Body.UserId),
	}

	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to create task")
	}

	id := uint(createdTask.ID)         // Преобразуем int64 в uint
	userId := uint(createdTask.UserID) // Преобразуем int64 в uint
	response := tasks.PostTasks201JSONResponse{
		Id:     &id,
		Task:   &createdTask.Task,
		IsDone: &createdTask.IsDone,
		UserId: &userId,
	}

	return response, nil
}

// PatchTasksId — обновление задачи по ID
func (h *Handler) PatchTasksId(_ context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	id := uint(request.Id)
	updatedTask := tasksService.Task{
		Task:   request.Body.Task,
		IsDone: request.Body.IsDone, // Используем значение напрямую
	}

	task, err := h.Service.UpdateTaskByID(id, updatedTask)
	if err != nil {
		return nil, err
	}

	taskID := uint(task.ID) // Преобразуем int64 в uint
	response := tasks.PatchTasksId200JSONResponse{
		Id:     &taskID,
		Task:   &task.Task,
		IsDone: &task.IsDone,
	}
	return response, nil
}

// DeleteTasksId — удаление задачи по ID
func (h *Handler) DeleteTasksId(_ context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	id := uint(request.Id)

	_, err := h.Service.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(404, "Task not found")
		}
		return nil, echo.NewHTTPError(500, "Failed to retrieve task")
	}

	err = h.Service.DeleteTaskByID(id)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to delete task")
	}

	return tasks.DeleteTasksId204Response{}, nil
}
