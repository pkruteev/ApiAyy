package queries

import (
	"ApiAyy/app/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) RegisterUser(b *models.UserType) error {

	// Определите строку запроса, исключив поле USER_ID.
	query := "INSERT INTO users (first_name, patronymic_name, last_name, user_email, user_phone, password) VALUES ($1, $2, $3, $4, $5, $6)"

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.UserPhone, b.Password)
	if err != nil {
		// Верните только ошибку.
		return err
	}

	return nil
}

func (q *UserQueries) SetupMember(id uint) error {

	const User_Right = "member"

	query := "INSERT INTO rights (user_id, user_right) VALUES ($1, $2)"
	fmt.Println(id, User_Right)

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, id, User_Right)
	if err != nil {
		// Верните только ошибку.
		return err
	}

	return nil
}

func (q *UserQueries) SetupAdmin(id uint) error {

	const User_Right = "admin"

	query := "INSERT INTO rights (user_id, user_right) VALUES ($1, $2)"
	fmt.Println(id, User_Right)

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, id, User_Right)
	if err != nil {
		// Верните только ошибку.
		return err
	}

	return nil
}

func (q *UserQueries) GetUserByEmail(UserEmail string) (models.UserType, error) {
	// fmt.Println(UserEmail)

	// Определяем переменную для хранения пользователя.
	user := models.UserType{}

	// Определяем строку запроса.
	query := "SELECT * FROM users WHERE user_email = $1"

	// Используйте QueryRow для получения строки.

	err := q.Get(&user, query, UserEmail)
	if err != nil {
		// Возвращаем пустой объект и ошибку.
		return user, err
	}

	// Возвращаем результат запроса.
	return user, nil
}

func (q *UserQueries) GetUserRightByID(userID uint) (string, error) {
	// Определяем переменную для хранения прав пользователя.
	var userRight string

	// Определяем строку запроса.
	query := "SELECT user_right FROM rights WHERE user_id = $1"

	// Используйте QueryRow для получения строки.
	err := q.QueryRow(query, userID).Scan(&userRight) // Используйте Scan для извлечения значения

	// Проверка на наличие ошибки
	if err != nil {
		return "", err // Возвращаем пустую строку и ошибку
	}

	return userRight, nil // Возвращаем права пользователя и nil как ошибку
}
