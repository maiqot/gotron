package tasksService

import "gorm.io/gorm"

// структура задачи

type Task struct {
	gorm.Model
	UserID uint   `json:"user_id" gorm:"not null"`
	Task   string `gorm:"not null"`
	IsDone bool   `gorm:"default:false"`
}

type TaskResponse struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserID uint   `json:"user_id"` // Добавляем user_id в ответ
}
