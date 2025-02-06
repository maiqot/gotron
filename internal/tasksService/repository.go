package tasksService

import "gorm.io/gorm"

type TaskRepository interface {
	// CreateTask - Передаем в функцию task типа Task из orm.go
	// возвращаем созданный Task и ошибку
	CreateTask(task Task) (Task, error)
	// GetAllTasks - Возвращаем массив из всех задач в БД и ошибку
	GetAllTasks() ([]Task, error)
	// UpdateTaskByID - Передаем id и Task, возвращаем обновленный Task
	// и ошибку
	UpdateTaskByID(id uint, task Task) (Task, error)
	// DeleteTaskByID - Передаем id для удаления, возвращаем только ошибку
	DeleteTaskByID(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// (r *taskRepository) привязывает данную функцию к нашему репозиторию
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

	// Ищем задачу по ID
	result := r.db.First(&task, id)
	if result.Error != nil {
		return task, result.Error
	}

	// Обновляем поля
	if updatedTask.Task != "" {
		task.Task = updatedTask.Task
	}
	task.IsDone = updatedTask.IsDone

	// Сохраняем изменения
	r.db.Save(&task)

	return task, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task

	// Ищем  задачу по ID
	result := r.db.First(&task, id)
	if result.Error != nil {
		return result.Error
	}

	// Удаляем задачу
	r.db.Delete(&task)

	return nil
}
