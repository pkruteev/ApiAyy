package models

import (
	"database/sql"
	"fmt"
)

// UserType определяет структуру пользователя
type UserType struct {
	UserID          int      `json:"user_id"`
	FirstName       string   `json:"first_name"`
	PatronymicName  string   `json:"patronymic_name"`
	LastName        string   `json:"last_name"`
	Discriminator   int      `json:"discriminator"` // 1 - administrator, 2 - director, 3 - manager
	UserEmail       string   `json:"user_email"`
	UserPhone       string   `json:"user_phone"`
	UserCompany     []string `json:"user_company,omitempty"`
	Password        string   `json:"password,omitempty"`
	ConfirmPassword string   `json:"confirm_password,omitempty"`
}

// SaveUser сохраняет пользователя в базе данных
func (u *UserType) SaveUser(db *sql.DB) error {
	query := `
    INSERT INTO users (first_name, patronymic_name, last_name, discriminator, user_email, user_phone, password)
    VALUES (?, ?, ?, ?, ?, ?, ?)
  `
	// Здесь вы можете использовать библиотеку для хеширования пароля, например, bcrypt
	hashedPassword := u.Password // Замените на хеширование пароля

	_, err := db.Exec(query, u.FirstName, u.PatronymicName, u.LastName, u.Discriminator, u.UserEmail, u.UserPhone, hashedPassword)
	if err != nil {
		return fmt.Errorf("could not save user: %v", err)
	}

	return nil
}
