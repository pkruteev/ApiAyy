package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetRSchets возвращает список всех расчетных счетов
func GetRSchets(c *fiber.Ctx) error {
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

	// Получаем расчетные счета
	rschets, err := db.GetRSchets()
	if err != nil {
		log.Printf("Ошибка при получении расчетных счетов: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении расчетных счетов: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"data":    rschets,
	})
}

// AddRSchet добавляет новый расчетный счет
func AddRSchet(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// 2. Проверка прав администратора
	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к основной БД: " + err.Error(),
		})
	}

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
		models.RSchet
		BDUsed string `json:"bd_used"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// 4. Подключение к целевой БД
	targetDB, err := database.DBConnectionQueries(request.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к целевой БД: " + err.Error(),
		})
	}

	// 5. Создание расчетного счета
	if err := targetDB.CreateRSchet(&request.RSchet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка создания расчетного счета: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"data":    request.RSchet,
	})
}

// DeleteRSchets удаляет расчетные счета по списку ID
func DeleteRSchets(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// 2. Проверка прав администратора
	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к основной БД: " + err.Error(),
		})
	}

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

	// 3. Парсинг входных данных
	var request struct {
		IDs []uint `json:"Ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Некорректный формат данных: " + err.Error(),
		})
	}

	// 4. Валидация входных данных
	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Массив ID не может быть пустым",
		})
	}

	// 5. Подключение к БД администратора
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	// 6. Удаление расчетных счетов
	if err := db.DeleteRSchets(request.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка удаления расчетных счетов: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"msg":     "Расчетные счета успешно удалены",
	})
}

// UpdateRSchet обновляет данные расчетного счета
func UpdateRSchet(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// 2. Проверка прав администратора
	dbMain, err := database.DBConnectionQueries("main")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к основной БД: " + err.Error(),
		})
	}

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
		models.RSchet
		BDUsed string `json:"bd_used"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// 4. Получение ID из URL
	rSchetIDStr := c.Params("id")
	rSchetID, err := strconv.ParseUint(rSchetIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный ID расчетного счета: " + rSchetIDStr,
		})
	}

	// 5. Подключение к БД администратора
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	// 6. Обновление расчетного счета
	request.RSchet.RSchetId = uint(rSchetID)
	if err := db.UpdateRSchet(request.RSchet.RSchetId, &request.RSchet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка обновления расчетного счета: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Расчетный счет обновлен успешно",
		"data":  request.RSchet,
	})
}
