package queries

import (
	"ApiAyy/app/models"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) RegisterUser(b *models.UserType) error {
	// Define query string.
	// query := `INSERT INTO users VALUES ( $2, $3, $4, $5, $6, $7, $7, $9)`
	// Определите строку запроса, исключив поле USER_ID.
	// Определите строку запроса, исключив поле USER_ID.
	query := "INSERT INTO users (first_name, patronymic_name, last_name, user_email, user_phone, user_company, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.UserPhone, b.UserCompany, b.Password)
	if err != nil {
		// Верните только ошибку.
		return err
	}

	return nil
}
