package userService

import (
	"firstProject/internal/tasksService"
	"time"
)

// структура пользователя

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"password" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Tasks []tasksService.Task `json:"tasks,omitempty" gorm:"foreignKey:UserID"`
}
