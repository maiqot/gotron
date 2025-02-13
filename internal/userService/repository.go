package userService

import "gorm.io/gorm"

// репозиторий пользователей

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUserByID(id uint, updatedUser User) (User, error) {
	var user User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}

	// Обновляем только измененные поля
	changes := map[string]interface{}{}
	if updatedUser.Email != "" {
		changes["email"] = updatedUser.Email
	}
	if updatedUser.Password != "" {
		changes["password"] = updatedUser.Password
	}

	// Если изменений нет — ничего не делаем
	if len(changes) == 0 {
		return user, nil
	}

	r.db.Model(&user).Updates(changes)
	return user, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}
