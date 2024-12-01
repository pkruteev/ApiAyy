package routes

import (
	"ApiAyy/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api")

	// Routes for GET method:

	route.Post("/register", controllers.Register)

	route.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": "AYY api",
			"version": "1.0.0",
		})
	})
}
