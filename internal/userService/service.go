package userService

import "firstProject/internal/tasksService"

type UserService struct {
	repo     UserRepository
	taskRepo tasksService.TaskRepository // taskRepo должен реализовывать TaskRepository
}

func NewUserService(repo UserRepository, taskRepo tasksService.TaskRepository) *UserService {
	return &UserService{repo: repo, taskRepo: taskRepo}
}

func (s *UserService) CreateUser(user User) (User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateUserByID(id uint, user User) (User, error) {
	return s.repo.UpdateUserByID(id, user)
}

func (s *UserService) DeleteUserByID(id uint) error {
	return s.repo.DeleteUserByID(id)
}

func (s *UserService) GetTasksForUser(userID uint) ([]tasksService.Task, error) {
	return s.taskRepo.GetTasksByUserID(userID) // Этот метод должен быть в TaskRepository
}
