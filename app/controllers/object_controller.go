package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetObjects возвращает список всех объектов из указанной базы данных.
func GetObjects(c *fiber.Ctx) error {
	// Получаем имя базы данных из параметра запроса
	bd := c.Params("bd")
	if bd == "" {
		log.Println("Имя базы данных не указано")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Имя базы данных не указано",
		})
	}

	// Подключаемся к указанной базе данных
	db, err := database.DBConnectionQueries(bd)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных %s: %v", bd, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к базе данных: " + err.Error(),
		})
	}

	// Получаем объекты из таблицы объектов
	objects, err := db.GetObjects()
	if err != nil {
		log.Printf("Ошибка при получении объектов: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении объектов: " + err.Error(),
		})
	}

	// Возвращаем список объектов в формате JSON
	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"data":    objects,
	})
}

// AddObject добавляет новый объект недвижимости в базу данных.
func AddObject(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к основной БД: " + err.Error(),
		})
	}

	// 2. Проверяем что пользователь - Администратор
	isAdmin, err := dbMain.IsAdmin(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка проверки прав: " + err.Error(),
		})
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "Недостаточно прав: требуется роль администратора",
		})
	}

	// 3. Парсинг входящих данных
	var request struct {
		models.Objects
		BDUsed string `json:"bd_used"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	targetDB, err := database.DBConnectionQueries(request.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	if err := targetDB.CreateObject(&request.Objects); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка создания объекта: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"msg":     "Объект успешно добавлен",
	})
}

// DeleteObjects удаляет объекты по списку ID
func DeleteObjects(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к основной БД: " + err.Error(),
		})
	}

	// 2. Проверяем что пользователь - Администратор
	isAdmin, err := dbMain.IsAdmin(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка проверки прав: " + err.Error(),
		})
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "Недостаточно прав: требуется роль администратора",
		})
	}

	// Парсинг входных данных
	var request struct {
		IDs []uint `json:"Ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Некорректный формат данных: " + err.Error(),
		})
	}

	// 3. Валидация входных данных
	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Массив Ids не может быть пустым",
		})
	}

	for _, id := range request.IDs {
		if id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "ID объекта не может быть нулевым",
			})
		}
	}

	// Подключение к БД администратора
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	// Удаление объектов
	failedDeletions := make([]uint, 0)
	successCount := 0

	for _, objectID := range request.IDs {
		// Проверка существования объекта
		exists, err := db.ObjectExists(objectID)
		if err != nil {
			log.Printf("Ошибка проверки существования объекта %d: %v", objectID, err)
			failedDeletions = append(failedDeletions, objectID)
			continue
		}

		if !exists {
			log.Printf("Объект с ID %d не найден", objectID)
			failedDeletions = append(failedDeletions, objectID)
			continue
		}

		// Удаление объекта
		if err := db.DeleteObject(objectID); err != nil {
			log.Printf("Ошибка удаления объекта %d: %v", objectID, err)
			failedDeletions = append(failedDeletions, objectID)
			continue
		}

		successCount++
	}

	// Формирование ответа
	response := fiber.Map{
		"success":      true,
		"deletedCount": successCount,
		"total":        len(request.IDs),
	}

	if len(failedDeletions) > 0 {
		response["failed"] = failedDeletions
		response["msg"] = fmt.Sprintf("Удалено %d из %d объектов", successCount, len(request.IDs))
	} else {
		response["msg"] = "Все объекты успешно удалены"
	}

	if successCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  true,
			"msg":    "Не удалось удалить ни одного объекта",
			"failed": failedDeletions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
