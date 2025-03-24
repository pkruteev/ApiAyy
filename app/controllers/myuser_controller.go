package controllers

import (
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Получить всех пользователей из одной компании
// Получить всех пользователей из одной компании
func GetMyUsers(c *fiber.Ctx) error {
	// Получаем текущее время.
	now := time.Now().Unix()

	// Получаем данные из JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Возвращаем статус 500 и ошибку парсинга JWT.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Устанавливаем время истечения токена.
	expires := claims.Expires
	userId := uint(claims.UserId)

	// Проверяем, истек ли токен.
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Подключение к базе данных main.
	bd_main, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем всех пользователей из таблицы rights, у которых user_bd равно bd_used.
	my_users, err := bd_main.GetMyUsers(userId)
	if err != nil {
		// Возвращаем статус 404, если пользователи не найдены.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "ошибка GetMyUsers, таблицы rights: " + err.Error(),
		})
	}

	// Возвращаем ответ
	return c.JSON(fiber.Map{
		"error": false,
		"users": my_users,
	})
}

// Записать в БД права существующего пользователя
func AddMyuser(c *fiber.Ctx) error {
	// Получаем текущее время.
	now := time.Now().Unix()

	// Получаем данные из JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	expires := claims.Expires

	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	userId := uint(claims.UserId)

	type Request struct {
		UserEmail string `db:"user_email"  json:"userEmail"`
		UserRole  string `db:"user_role"   json:"userRole,omitempty"`
	}
	myUserData := &Request{}
	if err := c.BodyParser(myUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Проверяем, что роль является допустимой
	validRoles := map[string]bool{
		"admin":    true,
		"member":   true,
		"director": true,
		"manager":  true,
	}
	if !validRoles[myUserData.UserRole] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("недопустимая роль: %s", myUserData.UserRole),
		})
	}
	// Подключение к базе данных main.
	db_main, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// получаем пользователя по EMAIL
	my_user, err := db_main.GetUserByEmail(myUserData.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось получить пользователя по EMAIL",
		})
	}

	// Записываем права в таблицу rights.
	err = db_main.SetupUserRight(my_user.UserID, userId, myUserData.UserRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось записать права пользователя в БД",
		})
	}

	// Возвращаем ответ
	return c.JSON(fiber.Map{
		"error":  false,
		"result": "ok",
	})
}
