package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetCounterparties возвращает список всех контрагентов
func GetCounterparties(c *fiber.Ctx) error {
	bd := c.Params("bd")
	if bd == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Имя базы данных не указано",
		})
	}

	db, err := database.DBConnectionQueries(bd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	counterparties, err := db.GetCounterparties()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка получения контрагентов: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    counterparties,
	})
}

// GetCounterparty возвращает одного контрагента по ID
func GetCounterparty(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный ID контрагента",
		})
	}

	bd := c.Params("bd")
	if bd == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Имя базы данных не указано",
		})
	}

	db, err := database.DBConnectionQueries(bd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	counterparty, err := db.GetCounterparty(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Контрагент не найден",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    counterparty,
	})
}

// Добавляет нового контрагента в бд.
type CreateCounterpartyRequest struct {
	BDUsed  string         `json:"bd_used"` // только дополнительные поля
	Company models.Company `json:"company"` // основная модель
}

func AddCounterparty(c *fiber.Ctx) error {

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

	// Проверяем что пользователь - Администратор
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
		models.Company
		BDUsed string `json:"bd_used"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных",
		})
	}

	// Подключаемся к целевой БД
	targetDB, err := database.DBConnectionQueries(request.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к целевой БД: " + err.Error(),
		})
	}

	if err := targetDB.CreateCounterparty(&request.Company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка создания контрагента",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    request,
	})
}

// UpdateCounterparty обновляет данные контрагента
func UpdateCounterparty(c *fiber.Ctx) error {

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
	// Проверяем что пользователь - Администратор
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
		models.Company
		BDUsed string `json:"bd_used"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// 4. Получение ID компании из URL
	companyIDStr := c.Params("id")
	fmt.Printf("Raw company ID from URL: '%s'\n", companyIDStr) // Логируем

	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный ID компании: " + companyIDStr,
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

	// 6. Обновление компании
	request.Company.CompanyId = uint(companyID)
	if err := db.UpdateCounterparty(request.Company.CompanyId, &request.Company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка обновления контрагента: " + err.Error(),
		})
	}

	// 7. Формирование ответа
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Компания обновлена успешно",
		"data":  request.Company,
	})
}

// DeleteCounterparty удаляет контрагента по списку ID
func DeleteCounterparties(c *fiber.Ctx) error {

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
			"msg":   "Массив companyIds не может быть пустым",
		})
	}

	for _, id := range request.IDs {
		if id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "ID контрагентов не может быть нулевым",
			})
		}
	}

	// Подключение к БД администратора
	userIdStr := strconv.FormatUint(uint64(userId), 10) // Преобразование uint в string
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	// Удаление контрагентов
	failedDeletions := make([]uint, 0)
	successCount := 0

	for _, companyID := range request.IDs {
		// Проверка существования контрагента
		exists, err := db.CounterpartyExists(companyID)
		if err != nil {
			log.Printf("Ошибка проверки существования контрагента %d: %v", companyID, err)
			failedDeletions = append(failedDeletions, companyID)
			continue
		}

		if !exists {
			log.Printf("Контрагент с ID %d не найдена", companyID)
			failedDeletions = append(failedDeletions, companyID)
			continue
		}

		// Удаление контрагента
		if err := db.DeleteCounterparty(companyID); err != nil {
			log.Printf("Ошибка удаления контрагента %d: %v", companyID, err)
			failedDeletions = append(failedDeletions, companyID)
			continue
		}

		successCount++
	}

	// 6. Формирование ответа
	response := fiber.Map{
		"success":      true,
		"deletedCount": successCount,
		"total":        len(request.IDs),
	}

	if len(failedDeletions) > 0 {
		response["failed"] = failedDeletions
		response["msg"] = fmt.Sprintf("Удалено %d из %d компаний", successCount, len(request.IDs))
	} else {
		response["msg"] = "Все компании успешно удалены"
	}

	if successCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  true,
			"msg":    "Не удалось удалить ни одной компании",
			"failed": failedDeletions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
