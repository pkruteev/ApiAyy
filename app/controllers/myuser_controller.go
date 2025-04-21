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

// Получить всех пользователей одной компании
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
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем всех пользователей из таблицы rights, у которых user_bd равно userId.
	userIdString := strconv.FormatUint(uint64(userId), 10)
	my_users, err := bd_main.GetMyAllUsers(userIdString)
	if err != nil {
		// Возвращаем статус 200 с сообщением, что пользователи не найдены
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"msg":   "Пока нет пользователей в данной базе данных",
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
	// Получаем текущее время
	now := time.Now().Unix()

	// Получаем данные из JWT
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

	// Уже получаем userId как uint, преобразование не нужно
	userId := uint(claims.UserId)
	userIdString := strconv.FormatUint(uint64(userId), 10) // Преобразуем uint в string

	type Request struct {
		UserEmail string `json:"userEmail"`
		UserRole  string `json:"userRole,omitempty"`
	}

	myUserData := &Request{}
	if err := c.BodyParser(myUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Подключение к базе данных main
	db_main, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем пользователя по EMAIL
	my_user, err := db_main.GetUserByEmail(myUserData.UserEmail)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось получить пользователя по EMAIL: " + err.Error(),
		})
	}

	// Записываем права в таблицу rights
	err = db_main.SetupUserRight(my_user.UserID, userIdString, myUserData.UserRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "не удалось записать права пользователя в БД: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":  false,
		"result": "ok",
	})
}

func DeleteMyUsers(c *fiber.Ctx) error {
	// 1. Проверка аутентификации
	if _, err := ValidateToken(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Требуется авторизация",
		})
	}

	// 2. Парсинг входных данных (ожидаем массив rightsId)
	var rightsIDs []uint
	if err := c.BodyParser(&rightsIDs); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Некорректный формат данных. Ожидается массив rightsIds",
		})
	}

	// 3. Валидация входных данных
	if len(rightsIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Массив rightsIds не может быть пустым",
		})
	}

	for _, id := range rightsIDs {
		if id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "ID не может быть нулевым",
			})
		}
	}

	// 4. Подключение к БД
	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Внутренняя ошибка сервера",
		})
	}
	defer dbMain.Close()

	// 5. Удаление записей
	deletedCount := 0
	for _, rightsID := range rightsIDs {
		exists, err := dbMain.RightsExists(rightsID)
		if err != nil {
			log.Printf("Ошибка проверки прав для rightsID %d: %v", rightsID, err)
			continue // Пропускаем проблемные записи
		}

		if !exists {
			log.Printf("Запись rightsID %d не найдена", rightsID)
			continue
		}

		if err := dbMain.DeleteRightsMyUser(rightsID); err != nil {
			log.Printf("Ошибка удаления rightsID %d: %v", rightsID, err)
			continue
		}

		deletedCount++
		log.Printf("Успешно удален rightsID: %d", rightsID)
	}

	if deletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Не удалось удалить ни одной записи",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"msg":     fmt.Sprintf("Удалено %d из %d записей", deletedCount, len(rightsIDs)),
	})
}
