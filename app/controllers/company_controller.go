package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetCompanies возвращает список всех компаний (counterparty: false)
func GetCompanies(c *fiber.Ctx) error {
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

	// Получаем компании из таблицы companies
	companies, err := db.GetCompanies()
	if err != nil {
		log.Printf("Ошибка при получении компаний: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка при получении компаний: " + err.Error(),
		})
	}

	// Возвращаем список компаний в формате JSON (пустой или нет)
	return c.JSON(fiber.Map{
		"success": true,
		"error":   false,
		"data":    companies,
	})
}

// AddCompany добавляет новую компанию в базу данных.
type CreateCompanyRequest struct {
	BDUsed  string         `json:"bd_used"` // только дополнительные поля
	Company models.Company `json:"company"` // основная модель
}

func AddCompany(c *fiber.Ctx) error {
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
		models.Company
		BDUsed string `json:"bd_used"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	// 4. Записываем имя компании в holding (из прав пользователя)
	rights, err := dbMain.GetUserRights(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка получения прав пользователя: " + err.Error(),
		})
	}

	if rights == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "Права пользователя не найдены",
		})
	}

	if rights.Holding == "" {
		// Если holding не указан, используем название компании
		rights.Holding = request.Company.Name
		if err := dbMain.UpdateUserHolding(userId, rights.Holding); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "Ошибка обновления holding: " + err.Error(),
			})
		}
	}

	targetDB, err := database.DBConnectionQueries(request.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к целевой БД: " + err.Error(),
		})
	}

	if err := targetDB.CreateCompany(&request.Company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка создания компании: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"holding": rights.Holding,
		},
		"error": false,
	})
}

// DeleteCompany удаляет компании по списку ID
func DeleteCompanies(c *fiber.Ctx) error {

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
				"msg":   "ID компании не может быть нулевым",
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

	// Удаление компаний
	failedDeletions := make([]uint, 0)
	successCount := 0

	for _, companyID := range request.IDs {
		// Проверка существования компании
		exists, err := db.CompanyExists(companyID)
		if err != nil {
			log.Printf("Ошибка проверки существования компании %d: %v", companyID, err)
			failedDeletions = append(failedDeletions, companyID)
			continue
		}

		if !exists {
			log.Printf("Компания с ID %d не найдена", companyID)
			failedDeletions = append(failedDeletions, companyID)
			continue
		}

		// Удаление компании
		if err := db.DeleteCompany(companyID); err != nil {
			log.Printf("Ошибка удаления компании %d: %v", companyID, err)
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

// Обновление по ID
func UpdateCompany(c *fiber.Ctx) error {
	// 1. Аутентификация и проверка прав
	userId, err := ValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// 2. Проверяем что пользователь - Администратор
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
	if err := db.UpdateCompany(request.Company.CompanyId, &request.Company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка обновления компании: " + err.Error(),
		})
	}

	// 7. Формирование ответа
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Компания обновлена успешно",
		"data":  request.Company,
	})
}
