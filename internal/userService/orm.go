package userService

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
