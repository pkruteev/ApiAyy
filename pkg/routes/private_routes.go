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
	route.Post("/rschets", middleware.JWTProtected(), controllers.AddRSchet)
	route.Post("/counterparties", middleware.JWTProtected(), controllers.AddCounterparty)
	route.Post("/addobject", middleware.JWTProtected(), controllers.AddObject)
	route.Post("/contracts", middleware.JWTProtected(), controllers.AddContract)

	// Routes for GET method:
	route.Get("/myusers", middleware.JWTProtected(), controllers.GetMyUsers)
	route.Get("/companies/:bd", middleware.JWTProtected(), controllers.GetCompanies)
	route.Get("/rschets/:bd", middleware.JWTProtected(), controllers.GetRSchets)
	route.Get("/counterparties/:bd", middleware.JWTProtected(), controllers.GetCounterparties)
	route.Get("/objects/:bd", middleware.JWTProtected(), controllers.GetObjects)
	route.Get("/contracts/:bd", middleware.JWTProtected(), controllers.GetContracts)
	// route.Get("/companies", middleware.JWTProtected(), controllers.GetCompanies)
	//route.Get("/companies/:bd", middleware.JWTProtected(), controllers.GetCompanies)

	// // Routes for PUT method:
	// route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID
	route.Put("/companies/:id", middleware.JWTProtected(), controllers.UpdateCompany)
	route.Put("/rschets/:id", middleware.JWTProtected(), controllers.UpdateRSchet)
	route.Put("/counterparties/:id", middleware.JWTProtected(), controllers.UpdateCounterparty)
	route.Put("/object/:id", middleware.JWTProtected(), controllers.UpdateObject)
	route.Put("/contracts/:id", middleware.JWTProtected(), controllers.UpdateContract)

	// // Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook)
	route.Delete("/companies", middleware.JWTProtected(), controllers.DeleteCompanies)
	route.Delete("/rschets", middleware.JWTProtected(), controllers.DeleteRSchets)
	route.Delete("/counterparties", middleware.JWTProtected(), controllers.DeleteCounterparties)
	route.Delete("/myusers", middleware.JWTProtected(), controllers.DeleteMyUsers)
	route.Delete("/objects", middleware.JWTProtected(), controllers.DeleteObjects)
	route.Delete("/contracts", middleware.JWTProtected(), controllers.DeleteContracts)
}
