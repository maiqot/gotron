package tasksService

import "errors"

// логика работы с задачами

type TaskService struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task Task) (Task, error) {
	// Проверяем, указан ли user_id
	if task.UserID == 0 {
		return Task{}, errors.New("user_id is required")
	}
	return s.repo.CreateTask(task)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTasksByUserID(userID uint) ([]Task, error) {
	return s.repo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTaskByID(id uint, task Task) (Task, error) {
	return s.repo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTaskByID(id uint) error {
	return s.repo.DeleteTaskByID(id)
}

// Новый метод для получения задачи по ID
func (s *TaskService) GetTaskByID(id uint) (Task, error) {
	return s.repo.GetTaskByID(id)
}
