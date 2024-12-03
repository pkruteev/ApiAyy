package controllers

import (
	"ApiAyy/pkg/utils"
)

func GetNewAccessToken(userID uint) (string, error) {
	// Генерируем новый Access токен, используя user ID.
	token, err := utils.GenerateNewAccessToken(userID) // Передаем userID в функцию генерации токена
	if err != nil {
		return "", err // Возвращаем ошибку
	}

	return token, nil // Возвращаем сгенерированный токен
}
