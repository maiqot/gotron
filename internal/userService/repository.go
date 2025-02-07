package userService

import "gorm.io/gorm"

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

	user.Email = updatedUser.Email
	user.Password = updatedUser.Password

	r.db.Save(&user)
	return user, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
