package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5" // Убедитесь, что вы используете актуальную версию пакета
)

// GenerateNewAccessToken генерирует новый Access токен для конкретного пользователя.
func GenerateNewAccessToken(userId uint) (string, error) {
	// Получаем секретный ключ из переменных окружения
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("секретный ключ не установлен")
	}

	// Получаем количество минут, на которое токен будет действовать, из переменных окружения
	minutesCountStr := os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	if minutesCountStr == "" {
		return "", fmt.Errorf("время жизни токена не установлено")
	}

	minutesCount, err := strconv.Atoi(minutesCountStr)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования времени жизни токена: %v", err)
	}

	// Создаем новые утверждения (claims) для токена
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix(),
		"user_id": userId,
	}

	// Создаем новый JWT токен с утверждениями
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Генерируем токен
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена: %v", err)
	}

	return t, nil
}
