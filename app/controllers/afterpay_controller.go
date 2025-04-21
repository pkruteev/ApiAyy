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

	// Устанавливаем время истечения токена.
	expires := claims.Expires
	userId := fmt.Sprintf("%v", claims.UserId)

	// Проверяем, истек ли токен.
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Подключение к основной базе данных.
	db_main, err := database.DBConnection("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Создание базы данных, если она еще не существует.
	err = database.CreateDB(userId, db_main)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка сохранения пользователя в БД: " + err.Error(),
		})
	}

	// Подключение к новой базе данных.
	db_new, err := database.DBConnection(userId)
	if err != nil {
		log.Fatalf("Ошибка подключения к новой базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к новой базе данных: " + err.Error(),
		})
	}

	// Создание таблиц
	if err := database.CreateTables(db_new); err != nil {
		log.Fatal("Ошибка создания таблиц:", err)
	}

	log.Println("Таблицы успешно созданы!")

	// Приведение типов для userId
	// userIdUint, err := strconv.ParseUint(userId, 10, 32)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "Некорректный userId: " + err.Error(),
	// 	})
	// }

	// Преобразуем строку userId в uint
	userIdUint64, err := strconv.ParseUint(userId, 10, 32) // 32 бита для соответствия uint
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Некорректный ID пользователя",
		})
	}

	// Явное преобразование в uint (32-битное беззнаковое целое)
	userIdUint := uint(userIdUint64)

	// Если нужно гарантированно получить строку без лишних символов:
	userIdString := strconv.FormatUint(userIdUint64, 10)

	// Преобразование в uint
	// userIdUint := uint(userIdUint)

	// Подключение к основной базе данных.
	db_main_queries, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Получаем права пользователя из таблицы Rights
	userRights, err := db_main_queries.GetUserRightsById(userIdUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении прав пользователя: " + err.Error(),
		})
	}
	log.Printf("userRights: %+v\n", userRights)
	// Проверяем, что пользователь имеет права admin и только в одной базе данных
	hasAdminRights := false
	for _, right := range userRights {
		if right.UserRole == "admin" {
			hasAdminRights = true
			break
		}
	}

	// Если пользователь уже является администратором в любой базе данных
	if hasAdminRights {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Пользователь с ID " + userId + " уже является администратором в другой базе данных",
		})
	}

	// Запись в поле bd_used имени основной бд для admin.
	err = db_main_queries.SetupUserBd(userIdUint, userIdString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка записи имени основной бд для admin: " + err.Error(),
		})
	}

	// Записываем пользовательские права admin в БД.
	err = db_main_queries.SetupUserRight(userIdUint, userIdString, "admin")
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Пользователь с такими правами уже зарегистрирован в системе",
		})
	}

	// Получаем обновленный статус пользователя из БД.
	updatedUserRights, err := db_main_queries.GetUserRightsById(userIdUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Не удалось получить статус пользователя из БД после сохранения",
		})
	}

	// Формируем ответ.
	response := fiber.Map{
		"message": fmt.Sprintf("База данных '%s' успешно создана и подключена", userId),
		"rights":  updatedUserRights,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
