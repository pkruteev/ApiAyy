package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetCompanies возвращает список всех компаний из указанной базы данных.
func GetCompanies(c *fiber.Ctx) error {
	// Получаем имя базы данных из параметра запроса
	bd := c.Params("bd")
	if bd == "" {
		log.Println("Имя базы данных не указано")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Имя базы данных не указано",
		})
	}

	// Подключаемся к указанной базе данных
	db, err := database.DBConnectionQueries(bd)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных %s: %v", bd, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем компании из таблицы companies
	companies, err := db.GetCompanies()
	if err != nil {
		log.Printf("Ошибка при получении компаний: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении компаний: " + err.Error(),
		})
	}

	// Возвращаем список компаний в формате JSON (пустой или нет)
	return c.JSON(fiber.Map{
		"error":     false,
		"companies": companies,
	})
}

// AddCompany добавляет новую компанию в базу данных.
type CreateCompanyRequest struct {
	BDUsed  string         `json:"bd_used"` // только дополнительные поля
	Company models.Company `json:"company"` // основная модель
}

func AddCompany(c *fiber.Ctx) error {
	// Проверка токена
	if _, err := ValidateToken(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Проверяем что пользователь - Администратор

	// Парсим данные сразу в models.Company
	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// bd_used
	var requestData struct {
		BDUsed string `json:"bd_used"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат bd_used: " + err.Error(),
		})
	}
	// Проверяем наличие записи в Holding таблицы RIGHTS

	// Подключаемся к БД
	db, err := database.DBConnectionQueries(requestData.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	// Создаем компанию
	if err := db.CreateCompany(&company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при создании компании: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    company.CompanyId,
		"error":   false,
	})
}
