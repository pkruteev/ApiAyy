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

func ValidateToken(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return err
	}

	expires := claims.Expires

	if now > expires {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized, check expiration time of your token")
	}

	return nil
}
