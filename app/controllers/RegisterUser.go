package controllers

import (
	"ApiAyy/app/models" // Импортируйте модель пользователя

	"github.com/gofiber/fiber/v2"
)

// RegisterUser создает нового пользователя
func RegisterUser(c *fiber.Ctx) error {
	var user models.UserType

	// Парсинг JSON тела запроса в структуру user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Здесь вы можете добавить логику валидации, проверки и сохранения пользователя
	// Пример: проверка, что пароли совпадают
	if user.Password != user.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	// Логика для сохранения пользователя в базу данных
	// Если сохранение прошло успешно, возвращаем успешный ответ
	// Пример:
	// err := saveUserToDatabase(user)
	// if err != nil {
	//   return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//     "error": "Could not register user",
	//   })
	// }

	// Успешный ответ
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}
