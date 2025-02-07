package handlers

import (
	"context"
	"firstProject/internal/userService"
	"firstProject/internal/web/users"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *userService.UserService
}

// Получение всех пользователей
func (h *UserHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := users.GetUsers200JSONResponse{}
	for _, user := range allUsers {
		user := users.User{
			Id:    &user.ID,
			Email: &user.Email,
		}
		response = append(response, user)
	}

	return response, nil
}

// Создание нового пользователя
func (h *UserHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	if request.Body == nil || request.Body.Email == nil || request.Body.Password == nil {
		return nil, echo.NewHTTPError(400, "Email and Password required")
	}

	userToCreate := userService.User{
		Email:    *request.Body.Email,
		Password: *request.Body.Password,
	}

	createdUser, err := h.Service.CreateUser(userToCreate)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to create user")
	}

	response := users.PostUsers201JSONResponse{
		Id:    &createdUser.ID,
		Email: &createdUser.Email,
	}
	return response, nil
}

// Обновление пользователя по ID
func (h *UserHandler) PatchUsersId(_ context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	id := uint(request.Id)

	if request.Body == nil || (request.Body.Email == nil && request.Body.Password == nil) {
		return nil, echo.NewHTTPError(400, "At least one field (Email or Password) is required")
	}

	updatedUser := userService.User{}
	if request.Body.Email != nil {
		updatedUser.Email = *request.Body.Email
	}
	if request.Body.Password != nil {
		updatedUser.Password = *request.Body.Password
	}

	user, err := h.Service.UpdateUserByID(id, updatedUser)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to update user")
	}

	response := users.PatchUsersId200JSONResponse{
		Id:    &user.ID,
		Email: &user.Email,
	}
	return response, nil
}

// Удаление пользователя по ID
func (h *UserHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	id := uint(request.Id)

	err := h.Service.DeleteUserByID(id)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to delete user")
	}

	return users.DeleteUsersId204Response{}, nil
}

// Конструктор хендлеров
func NewUserHandlers(service *userService.UserService) *UserHandler {
	return &UserHandler{Service: service}
}
