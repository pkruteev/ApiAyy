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

	// Создаем соединение с базой данных.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Сохраняем пользователя в базе данных
	err = db.RegisterUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка сохранения пользователя в БД: " + err.Error(),
		})
	}

	// Получаем пользователя из базы данных по email.
	createdUser, err := db.GetUserByUserEmail(user.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Неправильный логин или пароль",
		})
	}

	// Получаем идентификатор нового пользователя
	userID := user.User_ID

	// Генерация JWT токена
	token, err := GetNewAccessToken(userID) // Передаем userID в функцию
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка генерации токена: " + err.Error(),
		})
	}

	// Удаляем пароль из ответа
	createdUser.Password = ""
	fmt.Printf("Созданный пользователь: %+v\n", createdUser)

	// Возвращаем статус 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User registered successfully",
		"token": token,
		"user":  createdUser,
	})
}
