package routes

import (
	"ApiAyy/app/controllers"
	"ApiAyy/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api")

	// Routes for POST method:
	// route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)
	route.Post("/company", middleware.JWTProtected(), controllers.CreateCompany)
	route.Post("/afterpay", middleware.JWTProtected(), controllers.AfterPay)

	// // Routes for PUT method:
	// route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// // Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
}
