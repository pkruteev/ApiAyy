package utils

import (
	"ApiAyy/app/models"
	"ApiAyy/platform/database"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ConnectToUserDB function to connect to the user's database
func ConnectToUserDB(c *fiber.Ctx, user *models.UserType) (*database.Queries, error) {

	// Проверяем, валидны ли полученные данные JSON.
	if err := c.BodyParser(user); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userBd := user.UserBD

	// Проверяем, что userBd не равен 0
	if userBd == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Вы не авторизованы ни в одном холдинге")
	}

	// Преобразуем userBd в строку
	userBdStr := strconv.FormatUint(uint64(userBd), 10)

	// Подключение к базе данных пользователя.
	bd, err := database.DBConnectionQueries(userBdStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Ошибка подключения к базе данных: "+err.Error())
	}

	return bd, nil
}
