package handlers

import (
	"context"
	"firstProject/internal/userService"
	"firstProject/internal/web/users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/oapi-codegen/runtime/types"
)

type UserHandler struct {
	Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

// Получение всех пользователей
// GetUsers retrieves all users
func (h *UserHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var response users.GetUsers200JSONResponse
	for _, u := range allUsers {
		id := uint(u.ID) // Преобразуем int64 в uint
		user := users.User{
			Id:    &id,
			Email: (*types.Email)(&u.Email),
		}
		response = append(response, user)
	}

	return response, nil
}

// PostUsers creates a new user
func (h *UserHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	if request.Body == nil || request.Body.Email == "" || request.Body.Password == "" {
		return nil, echo.NewHTTPError(400, "Email and Password required")
	}

	userToCreate := userService.User{
		Email:    string(request.Body.Email), // Убираем лишнее разыменование
		Password: request.Body.Password,
	}

	createdUser, err := h.Service.CreateUser(userToCreate)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to create user")
	}

	id := uint(createdUser.ID) // Преобразуем int64 в uint
	response := users.PostUsers201JSONResponse{
		Id:    &id,
		Email: (*types.Email)(&createdUser.Email),
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
		updatedUser.Email = string(*request.Body.Email)
	}
	if request.Body.Password != nil {
		updatedUser.Password = *request.Body.Password
	}

	user, err := h.Service.UpdateUserByID(id, updatedUser)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to update user")
	}

	response := users.PatchUsersId200JSONResponse{
		Id:    &id,
		Email: (*types.Email)(&user.Email),
	}
	return response, nil
}

// Удаление пользователя по ID
func (h *UserHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	id := uint(request.Id)

	tasks, _ := h.Service.GetTasksForUser(id)
	log.Printf("Задачи пользователя %d: %+v", id, tasks)

	// Удаляем все задачи, связанные с пользователем
	err := h.Service.DeleteTasksByUserID(id)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to delete user's tasks")
	}

	// Удаляем пользователя
	err = h.Service.DeleteUserByID(id)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Failed to delete user")
	}

	return users.DeleteUsersId204Response{}, nil
}
