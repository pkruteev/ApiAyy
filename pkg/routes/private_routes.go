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
	// route.Post("/company", middleware.JWTProtected(), controllers.CreateCompany)
	route.Post("/afterpay", middleware.JWTProtected(), controllers.AfterPay)
	route.Post("/addmyuser", middleware.JWTProtected(), controllers.AddMyuser)
	route.Post("/addcompany", middleware.JWTProtected(), controllers.AddCompany)

	// Routes for GET method:
	route.Get("/myusers", middleware.JWTProtected(), controllers.GetMyUsers)
	route.Get("/companies/:bd", middleware.JWTProtected(), controllers.GetCompanies)
	route.Get("/objects/:bd", middleware.JWTProtected(), controllers.GetObjects)
	// route.Get("/companies", middleware.JWTProtected(), controllers.GetCompanies)
	//route.Get("/companies/:bd", middleware.JWTProtected(), controllers.GetCompanies)

	// // Routes for PUT method:
	// route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// // Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook)
	route.Delete("/companies", middleware.JWTProtected(), controllers.DeleteCompanies)
	route.Delete("/myusers", middleware.JWTProtected(), controllers.DeleteMyUsers)
}
