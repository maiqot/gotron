package tasksService

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Task   string `gorm:"not null"`
	IsDone bool   `gorm:"default:false"`
	UserID int    `json:"user_id" gorm:"not null"` // Добавляем связь с пользователем
}

type TaskResponse struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserID uint   `json:"user_id"` // Добавляем user_id в ответ
}
