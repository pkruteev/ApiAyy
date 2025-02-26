package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Login контроллер для входа пользователя в систему.
func LoginCopy(c *fiber.Ctx) error {
	// Создаем новую структуру для получения данных пользователя.
	loginData := &models.UserType{}
	if err := c.BodyParser(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	queries, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Получаем пользователя из базы данных по email.
	user, err := queries.GetUserByEmail(loginData.UserEmail)
	if err != nil {
		fmt.Println("Пользователь не найден по данному email:", loginData.UserEmail)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Неправильный логин или пароль",
		})
	}

	// Сравниваем введенный пароль с сохраненным хешем
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Неправильный логин или пароль",
		})
	}

	// Генерация JWT токена
	token, err := GetNewAccessToken(user.UserID) // Используйте идентификатор пользователя
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка генерации токена: " + err.Error(),
		})
	}

	// Удаляем пароль из ответа
	user.Password = ""
	// user.CreatedUser = ""

	// Возвращаем статус 200 OK, данные пользователя и токен
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Успешный вход в систему",
		"token": token,
		"user":  loginData,
	})
}
