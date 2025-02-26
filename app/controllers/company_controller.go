package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Получить все компании
func GetCompanies(c *fiber.Ctx) error {
	// Проверка токена
	if err := ValidateToken(c); err != nil {
		return err // Возвращаем ошибку, если токен недействителен или истек
	}

	user := &models.UserType{}
	// Проверяем, валидны ли полученные данные JSON.
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userBd := user.UserBD

	// Проверяем, что userBd не равен 0
	if userBd == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Вы не авторизованы ни в одном холдинге",
		})
	}

	// Преобразуем userBd в строку
	userBdStr := strconv.FormatUint(uint64(userBd), 10)

	// Подключение к базе данных пользователя.
	bd, err := database.DBConnectionQueries(userBdStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем компании из таблицы companies.
	companies, err := bd.GetCompanies()
	if err != nil {
		// Возвращаем статус 404, если компании не найдены.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Не удалось получить список компаний из БД",
		})
	}

	// Возвращаем ответ
	return c.JSON(fiber.Map{
		"error":     false,
		"companies": companies,
	})
}

// Записать в БД новую компанию
func AddCompany(c *fiber.Ctx) error {
	// Проверка токена
	if err := ValidateToken(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Недействительный или истекший токен",
		})
	}

	company := &models.Company{}
	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	user := &models.UserType{}
	// Подключаемся к базе данных
	db, err := utils.ConnectToUserDB(c, user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Вызовите метод CreateCompany
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
