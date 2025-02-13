package userService

import (
	"firstProject/internal/tasksService"
	"log"
)

// логика работы с пользователями
type UserService struct {
	repo     UserRepository
	taskRepo tasksService.TaskRepository
}

func NewUserService(repo UserRepository, taskRepo tasksService.TaskRepository) *UserService {
	return &UserService{repo: repo, taskRepo: taskRepo}
}

func (s *UserService) CreateUser(user User) (User, error) {
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
	}
	return createdUser, err
}

func (s *UserService) GetAllUsers() ([]User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		log.Printf("Ошибка при получении пользователей: %v", err)
	}
	return users, err
}

func (s *UserService) UpdateUserByID(id uint, user User) (User, error) {
	updatedUser, err := s.repo.UpdateUserByID(id, user)
	if err != nil {
		log.Printf("Ошибка при обновлении пользователя %d: %v", id, err)
	}
	return updatedUser, err
}

func (s *UserService) DeleteUserByID(id uint) error {
	err := s.repo.DeleteUserByID(id)
	if err != nil {
		log.Printf("Ошибка при удалении пользователя %d: %v", id, err)
	}
	return err
}

func (s *UserService) GetTasksForUser(userID uint) ([]tasksService.Task, error) {
	tasks, err := s.taskRepo.GetTasksByUserID(userID)
	if err != nil {
		log.Printf("Ошибка при получении задач пользователя %d: %v", userID, err)
	}
	return tasks, err
}

// Удаление всех задач пользователя
func (s *UserService) DeleteTasksByUserID(userID uint) error {
	tasks, err := s.taskRepo.GetTasksByUserID(userID)
	if err != nil {
		log.Printf("Ошибка при получении задач для удаления пользователя %d: %v", userID, err)
		return err
	}

	if len(tasks) == 0 {
		log.Printf("У пользователя %d нет задач для удаления", userID)
		return nil
	}

	// Просто удаляем задачи без транзакции
	for _, task := range tasks {
		if err := s.taskRepo.DeleteTaskByID(task.ID); err != nil {
			log.Printf("Ошибка при удалении задачи %d пользователя %d: %v", task.ID, userID, err)
			return err
		}
	}

	log.Printf("Все задачи пользователя %d удалены", userID)
	return nil
}
