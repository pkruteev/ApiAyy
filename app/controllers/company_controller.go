package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"log"
	"strconv"

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
func AddCompany(c *fiber.Ctx) error {
	// Проверка токена
	if _, err := ValidateToken(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	company := &models.Company{}

	// Структура для парсинга bd_used
	type RequestBody struct {
		BDUsed int `json:"bd_used"`
	}

	var body RequestBody

	// Парсим bd_used
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// Парсим company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// Преобразуем bd_used в строку
	bd := strconv.Itoa(body.BDUsed)

	// Подключаемся к базе данных
	db, err := database.DBConnectionQueries(bd)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных %s: %v", bd, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Вызываем метод CreateCompany
	err = db.CreateCompany(company)
	if err != nil {
		log.Printf("Ошибка при создании компании: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при создании компании: " + err.Error(),
		})
	}

	log.Println("Компания успешно создана!")

	return c.JSON(fiber.Map{
		"error":  false,
		"result": "Компания успешно создана!",
	})
}
