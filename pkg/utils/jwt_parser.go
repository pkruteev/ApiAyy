package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	Expires int64
	UserId  float64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// fmt.Println("Claims:", claims)

		// Извлечение user_id из claims
		userId := claims["user_id"].(float64)

		// Вывод значения user_id в терминал
		fmt.Printf("Extracted user_id: %.0f\n", userId)

		// Expires time.
		expires := int64(claims["exp"].(float64))

		return &TokenMetadata{
			Expires: expires,
			UserId:  userId, // Вы можете добавить это поле в TokenMetadata, если нужно
		}, nil
	}

	return nil, err
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
