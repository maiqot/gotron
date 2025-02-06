package handlers

import (
	"context"
	"firstProject/internal/tasksService" // Импортируем наш сервис
	"firstProject/internal/web/tasks"
)

type Handler struct {
	Service *tasksService.TaskService
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
	// Распаковываем тело запроса напрямую, без декодера!
	task := request.Body
	// Обращаемся к сервису и создаем задачу
	taskToCreate := tasksService.Task{
		Task:   *task.Task,
		IsDone: *task.IsDone,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)

	if err != nil {
		return nil, err
	}
	// создаем структуру респонс
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Task,
		IsDone: &createdTask.IsDone,
	}
	// Просто возвращаем респонс!
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

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *tasksService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}
