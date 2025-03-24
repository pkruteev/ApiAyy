package controllers

import (
	"ApiAyy/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetNewAccessToken(userID uint) (string, error) {
	// Генерируем новый Access токен, используя user ID.
	token, err := utils.GenerateNewAccessToken(userID) // Передаем userID в функцию генерации токена
	if err != nil {
		return "", err // Возвращаем ошибку
	}

	return token, nil // Возвращаем сгенерированный токен
}

func ValidateToken(c *fiber.Ctx) (uint, error) {
	// Получаем текущее время.
	now := time.Now().Unix()

	// Получаем данные из JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Возвращаем ошибку парсинга JWT.
		return 0, err
	}

	// Устанавливаем время истечения токена.
	expires := claims.Expires
	userId := uint(claims.UserId)

	// Проверяем, истек ли токен.
	if now > expires {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "unauthorized, check expiration time of your token")
	}

	return userId, nil
}
