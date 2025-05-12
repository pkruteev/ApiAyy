package controllers

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetContracts возвращает список всех договоров
func GetContracts(c *fiber.Ctx) error {
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

	contracts, err := db.GetContracts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка получения договоров: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    contracts,
	})
}

// GetContract возвращает один договор по ID
func GetContract(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный ID договора",
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

	contract, err := db.GetContract(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Договор не найден",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    contract,
	})
}

// AddContract добавляет новый договор
func AddContract(c *fiber.Ctx) error {
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

	var request struct {
		models.ContractType
		BDUsed string `json:"bd_used"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных",
		})
	}

	targetDB, err := database.DBConnectionQueries(request.BDUsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к целевой БД: " + err.Error(),
		})
	}

	if err := targetDB.CreateContract(&request.ContractType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка создания договора",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    request.ContractType,
	})
}

// UpdateContract обновляет данные договора
func UpdateContract(c *fiber.Ctx) error {
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

	var request struct {
		models.ContractType
		BDUsed string `json:"bd_used"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный формат данных: " + err.Error(),
		})
	}

	contractIDStr := c.Params("id")
	contractID, err := strconv.ParseUint(contractIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Неверный ID договора: " + contractIDStr,
		})
	}

	userIdStr := strconv.FormatUint(uint64(userId), 10)
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	request.ContractType.ContractId = uint(contractID)
	if err := db.UpdateContract(request.ContractType.ContractId, &request.ContractType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка обновления договора: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Договор обновлен успешно",
		"data":  request.ContractType,
	})
}

// DeleteContracts удаляет договоры по списку ID
func DeleteContracts(c *fiber.Ctx) error {
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

	var request struct {
		IDs []uint `json:"Ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Некорректный формат данных: " + err.Error(),
		})
	}

	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Массив contractIds не может быть пустым",
		})
	}

	for _, id := range request.IDs {
		if id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "ID договоров не может быть нулевым",
			})
		}
	}

	userIdStr := strconv.FormatUint(uint64(userId), 10)
	db, err := database.DBConnectionQueries(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Ошибка подключения к БД: " + err.Error(),
		})
	}

	failedDeletions := make([]uint, 0)
	successCount := 0

	for _, contractID := range request.IDs {
		exists, err := db.ContractExists(contractID)
		if err != nil {
			log.Printf("Ошибка проверки существования договора %d: %v", contractID, err)
			failedDeletions = append(failedDeletions, contractID)
			continue
		}

		if !exists {
			log.Printf("Договор с ID %d не найден", contractID)
			failedDeletions = append(failedDeletions, contractID)
			continue
		}

		if err := db.DeleteContract(contractID); err != nil {
			log.Printf("Ошибка удаления договора %d: %v", contractID, err)
			failedDeletions = append(failedDeletions, contractID)
			continue
		}

		successCount++
	}

	response := fiber.Map{
		"success":      true,
		"deletedCount": successCount,
		"total":        len(request.IDs),
	}

	if len(failedDeletions) > 0 {
		response["failed"] = failedDeletions
		response["msg"] = fmt.Sprintf("Удалено %d из %d договоров", successCount, len(request.IDs))
	} else {
		response["msg"] = "Все договоры успешно удалены"
	}

	if successCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  true,
			"msg":    "Не удалось удалить ни одного договора",
			"failed": failedDeletions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
