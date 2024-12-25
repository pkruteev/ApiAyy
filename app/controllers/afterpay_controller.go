package controllers

import (
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AfterPay(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current.
	expires := claims.Expires
	userId := fmt.Sprintf("%v", claims.UserId) // Приведение типа float64 к string
	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Подключение к основной базе данных
	db_main, err := database.PostgreSQLConnection()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Вызов функции для создания db с названием UsersID если не созданы
	database.CreateDB(userId, db_main)

	// Закрываем исходное соединение с базой данных db_main
	// err := db_main.Close()
	// if err != nil {
	// 	log.Printf("Ошибка при закрытии соединения с БД: %v", err)
	// 	return err
	// }

	// Закрываем исходное соединение с базой данных db_main после создания новой бд
	err = db_main.Close()
	if err != nil {
		log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		return err
	}

	// Подключение к новой базе данных PostgreNewSQLConnection

	// Подключение к новой базе данных
	db_new, err := database.PostgreNewSQLConnection(userId)
	if err != nil {
		log.Fatalf("Ошибка подключения к новой базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к новой базе данных: " + err.Error(),
		})
	}

	// Вызов функции для создания таблиц если не созданы
	database.CreateAllTables(db_new)

	////
	// db_new, err := database.PostgreNewSQLConnection(userId)
	// if err != nil {
	// 	log.Fatalf("Ошибка подключения к базе данных: %v", err)
	// }

	// Закрываем новое соединение с базой данных после работы с ней,
	// если необходимо.
	defer func() {
		if err := db_new.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с новой БД: %v", err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("База данных '%s' успешно создана и подключена", userId),
	})
}
