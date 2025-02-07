package main

import (
	"firstProject/internal/database"
	"firstProject/internal/handlers"
	"firstProject/internal/tasksService"
	"firstProject/internal/userService"
	"firstProject/internal/web/tasks"
	"firstProject/internal/web/users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	database.InitDB()

	// Миграция моделей
	if err := database.DB.AutoMigrate(&tasksService.Task{}, &userService.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Инициализация сервисов
	taskRepo := tasksService.NewTaskRepository(database.DB)
	taskService := tasksService.NewService(taskRepo)
	taskHandler := handlers.NewHandler(taskService)

	userRepo := userService.NewUserRepository(database.DB)
	userServiceInstance := userService.NewUserService(userRepo)
	userHandler := handlers.NewUserHandlers(userServiceInstance)

	// Инициализируем Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрируем обработчики
	strictTaskHandler := tasks.NewStrictHandler(taskHandler, nil)
	tasks.RegisterHandlers(e, strictTaskHandler)

	strictUserHandler := users.NewStrictHandler(userHandler, nil)
	users.RegisterHandlers(e, strictUserHandler)

	// Запуск сервера
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
