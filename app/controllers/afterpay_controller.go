package controllers

import (
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"strconv"
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
	////Блок создания таблиц по имени User ID/////
	// Подключение к базе данных main
	db, err := database.DBConnection("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	//defer db.Close()

	// Вызов функции для создания db с названием UsersID если не созданы
	database.CreateDB(userId, db)

	// Закрываем исходное соединение с базой данных db_main после создания новой бд
	err = db.Close()
	if err != nil {
		log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		return err
	}

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

	// Закрываем новое соединение с базой данных после работы с ней,

	defer func() {
		if err := db_new.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с новой БД: %v", err)
		}
	}()
	////Конец блока создания таблиц////
	//Создаем соединение с базой данных.
	db_main, err := database.DBMainConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Приведение типов.
	userIdUint64, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка преобразования userId в uint: " + err.Error(),
		})
	}

	userIdUint := uint(userIdUint64)

	// Записываем пользовательские права admin в БД.
	err = db_main.SetupUserRight(userIdUint, "admin", userIdUint)
	var setupErrorMsg string
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Err 409 - Пользователь с такими правами уже зарегистрирован в системе",
		})
	}

	// Получаем статус пользователя из БД
	rightsUser, err := db_main.GetUserRightByID(userIdUint)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Не удалось получить статус пользователя из БД после сохранения",
		})
	}

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message":    fmt.Sprintf("База данных '%s' успешно создана и подключена", userId),
	// 	"user_right": rightsUser,
	// })
	// Формируем ответ
	response := fiber.Map{
		"message":    fmt.Sprintf("База данных '%s' успешно создана и подключена", userId),
		"user_right": rightsUser,
	}

	// Если была ошибка при установке прав, добавляем ее в ответ
	if setupErrorMsg != "" {
		response["setup_error"] = setupErrorMsg
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
