package tasksService

import "gorm.io/gorm"

// репозиторий задач

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
	GetTasksByUserID(userID uint) ([]Task, error)
	GetTaskByID(id uint) (Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, updatedTask Task) (Task, error) {
	var task Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return task, result.Error
	}
	if updatedTask.Task != "" {
		task.Task = updatedTask.Task
	}
	task.IsDone = updatedTask.IsDone
	r.db.Save(&task)
	return task, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return result.Error
	}
	r.db.Delete(&task)
	return nil
}

func (r *taskRepository) GetTasksByUserID(userID uint) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTaskByID(id uint) (Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	return task, err
}
