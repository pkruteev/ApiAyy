package main

import (
	"ApiAyy/pkg/configs"
	"ApiAyy/pkg/middleware"
	"ApiAyy/pkg/routes"
	"ApiAyy/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app)

	// // Подключение к базе данных
	// db, err := database.OpenDBConnection()
	// if err != nil {
	// 	// Return status 500 and database connection error.
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   err.Error(),
	// 	})
	// }

	fmt.Println("Connection pool created successfully.")

	// Routes.
	// В случае, если вам нужно использовать объект db в маршрутах, передайте его в маршруты
	routes.PublicRoutes(app)  // Register a public routes for app, передаем db
	routes.PrivateRoutes(app) // Register a private routes for app, передаем db
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with graceful shutdown).
	utils.StartServer(app)
}
