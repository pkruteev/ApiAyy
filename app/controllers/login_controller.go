package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	// Создаем структуру для получения данных пользователя
	loginData := &models.UserType{}
	if err := c.BodyParser(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// Создаем соединение с базой данных main
	db_main, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем пользователя из базы данных по email
	user, err := db_main.GetUserByEmail(loginData.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении пользователя: " + err.Error(),
		})
	}

	// Если пользователь не найден (UserID == 0), возвращаем данные для регистрации
	if user.UserID == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"user":  user, // Возвращаем структуру с email для регистрации
		})
	}

	// Сравниваем введенный пароль с хешем из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Неправильный пароль",
		})
	}

	// Получаем права пользователя из таблицы Rights
	userRights, err := db_main.GetUserRightsById(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении прав пользователя: " + err.Error(),
		})
	}

	// Генерация JWT токена
	token, err := GetNewAccessToken(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка генерации токена: " + err.Error(),
		})
	}

	// Удаляем пароль из ответа
	user.Password = ""

	// Логируем вход пользователя (без чувствительных данных)
	fmt.Printf("Login user: %+v\n", user.UserEmail)

	// Возвращаем полные данные пользователя клиенту
	return c.JSON(fiber.Map{
		"error":  false,
		"token":  token,
		"user":   user,
		"rights": userRights,
	})
}
