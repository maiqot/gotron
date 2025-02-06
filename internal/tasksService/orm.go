package tasksService

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Task   string `gorm:"not null"`
	IsDone bool   `gorm:"default:false"`
}

type TaskResponse struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}
