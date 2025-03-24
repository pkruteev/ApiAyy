package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"

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
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// Шифруем пароль пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при шифровании пароля: " + err.Error(),
		})
	}

	// Устанавливаем зашифрованный пароль как строку
	user.Password = string(hashedPassword)

	// Создаем соединение с базой данных main.
	db_main, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}
	// defer db_main.Close() // Закрываем соединение после использования

	// Сохраняем пользователя в базе данных
	err = db_main.RegisterUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка сохранения пользователя в БД: " + err.Error(),
		})
	}

	// Получаем пользователя из базы данных по email.
	createdUser, err := db_main.GetUserByEmail(user.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении пользователя из БД: " + err.Error(),
		})
	}

	// Получаем идентификатор нового пользователя
	userID := createdUser.UserID

	// Генерация JWT токена
	token, err := GetNewAccessToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка генерации токена: " + err.Error(),
		})
	}

	// Удаляем пароль из ответа
	createdUser.Password = ""

	// Логируем регистрацию пользователя (без чувствительных данных)
	fmt.Printf("Зарегистрирован пользователь: %+v\n", createdUser.UserEmail)

	// Возвращаем полные данные пользователя клиенту
	return c.JSON(fiber.Map{
		"error": false,
		"token": token,
		"msg":   "Пользователь успешно зарегистрирован",
		"user":  createdUser,
	})
}
