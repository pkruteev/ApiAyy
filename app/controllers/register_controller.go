package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	// Создаем новую структуру пользователя

	user := &models.UserType{}
	// Проверяем, валидны ли полученные данные JSON.
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Шифруем пароль пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error processing password",
		})
	}

	// Устанавливаем зашифрованный пароль как строку
	user.Password = string(hashedPassword)

	//Создаем соединение с базой данных main.
	db_main, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Сохраняем пользователя в базе данных
	err = db_main.RegisterUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка сохранения пользователя в БД: " + err.Error(),
		})
	}

	// Получаем пользователя из базы данных по email.
	// createdUser, err := db_main.GetUserByEmail(user.UserEmail)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "Не удалось получить пользователя из БД",
	// 	})
	// }

	// Получаем идентификатор нового пользователя
	// userID := createdUser.UserID

	// Генерация JWT токена
	// token, err := GetNewAccessToken(userID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "Ошибка генерации токена: " + err.Error(),
	// 	})
	// }

	// Удаляем пароль из ответа
	// createdUser.Password = ""

	// fmt.Printf("Зарегистрирован пользователь: %+v\n", createdUser)

	// Возвращаем статус 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User registered successfully",
		// "token": token,
		// "user":  createdUser,
	})
}
