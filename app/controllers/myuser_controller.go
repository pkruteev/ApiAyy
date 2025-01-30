package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	myQueries, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем всех пользователей из таблицы rights.
	my_users, err := myQueries.GetMyUsers(userId)
	if err != nil {
		// Возвращаем статус 404, если пользователи не найдены.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось получить список пользователей из БД",
		})
	}

	// Массив для формируемых пользователей
	var userResponses []models.UserResponse

	// Получаем пользователей из двух таблиц для отправки
	for _, my_user := range my_users {
		userResponse, err := myQueries.GetUserForResponsById(my_user.User_ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "не удалось получить пользователя с ID " + fmt.Sprintf("%d", my_user.User_ID),
			})
		}

		// Проверяем, что данные возвращаемые из запроса совпадают с нужной структурой
		userResponseStruct := models.UserResponse{
			UserID:         userResponse.UserID,
			FirstName:      userResponse.FirstName,
			LastName:       userResponse.LastName,
			PatronymicName: userResponse.PatronymicName,
			UserRights:     userResponse.UserRights,
			UserEmail:      userResponse.UserEmail,
			UserPhone:      userResponse.UserPhone,
		}

		// Добавляем заполненный userResponse в массив
		userResponses = append(userResponses, userResponseStruct)
	}

	// Возвращаем ответ
	return c.JSON(fiber.Map{
		"error": false,
		"users": userResponses,
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

	myUserData := &models.UserType{}
	if err := c.BodyParser(myUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userId := uint(claims.UserId)

	// Подключение к базе данных main.
	myQueries, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// получаем пользователя по EMAIL
	my_user, err := myQueries.GetUserByEmail(myUserData.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось получить пользователя по EMAIL",
		})
	}

	// Записываем права в таблицу rights.
	err = myQueries.SetupUserRight(my_user.User_ID, userId, myUserData.User_Rights)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось записать права пользователя в БД",
		})
	}

	// Получаем пользователя из двух таблиц для отправки
	userResponse, err := myQueries.GetUserForResponsById(my_user.User_ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось получить пользователя",
		})
	}

	// Возвращаем ответ
	return c.JSON(fiber.Map{
		"error": false,
		"users": userResponse,
	})
}
