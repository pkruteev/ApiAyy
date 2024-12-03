package main

import (
	"ApiAyy/pkg/configs"
	"ApiAyy/pkg/middleware"
	"ApiAyy/pkg/routes"
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app)

	// Подключение к базе данных
	db, err := database.PostgreSQLConnection()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Вызов функции для создания таблиц если не созданы
	database.CreateTables(db)

	fmt.Println("Пул соединений успешно создан.")

	// Routes.
	// В случае, если вам нужно использовать объект db в маршрутах, передайте его в маршруты
	routes.PublicRoutes(app)  // Register a public routes for app, передаем db
	routes.PrivateRoutes(app) // Register a private routes for app, передаем db
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with graceful shutdown).
	utils.StartServer(app)
}
