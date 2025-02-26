package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) RegisterUser(b *models.UserType) error {

	// Строка запроса, без поля USER_ID.
	query := "INSERT INTO users ( bd_used, first_name, patronymic_name, last_name, user_email, user_phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := q.Exec(query, b.BdUsed, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.UserPhone, b.Password)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) SetupUserRight(id uint, userBd uint, rights string) error {

	query := "INSERT INTO rights (user_id, user_bd, holding, user_rights) VALUES ($1, $2, '', $3)"

	existsQuery := `
		 SELECT 
			  COUNT(*) FILTER (WHERE user_id = $1 AND user_bd = $2 AND user_rights = $3) AS rights_count,
			  COUNT(*) FILTER (WHERE user_id = $1 AND user_rights = 'admin') AS admin_count
		 FROM rights`

	var rightsCount, adminCount int
	err := q.QueryRow(existsQuery, id, userBd, rights).Scan(&rightsCount, &adminCount)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования записи: %w", err)
	}

	if rightsCount > 0 {
		return fmt.Errorf("пользователь с ID %d уже имеет необходимые права: %s для компании: %d", id, rights, userBd)
	}

	if adminCount > 0 {
		return fmt.Errorf("пользователь с ID %d уже имеет права администратора", id)
	}

	_, err = q.Exec(query, id, userBd, rights)
	if err != nil {
		return fmt.Errorf("ошибка при установке прав пользователя: %w", err)
	}

	return nil
}

func (q *UserQueries) GetUserByEmail(UserEmail string) (models.UserType, error) {

	// Определяем переменную для хранения пользователя.
	user := models.UserType{}
	// fmt.Printf("Email: %+v\n", UserEmail)
	// Определяем строку запроса.
	query := "SELECT * FROM users WHERE user_email = $1"

	// Используйте QueryRow для получения строки.
	err := q.Get(&user, query, UserEmail)
	if err != nil {
		// Возвращаем пустой объект и ошибку.
		fmt.Printf("Error retrieving user: %v\n", err)
		return user, err
	}

	// Возвращаем результат запроса.
	return user, nil
}

func (q *UserQueries) GetUserRightByID(userID uint) (string, error) {
	// Определяем переменную для хранения прав пользователя.
	var userRights string

	// Определяем строку запроса.
	query := "SELECT user_rights FROM rights WHERE user_id = $1"

	// Используйте QueryRow для получения строки.
	err := q.QueryRow(query, userID).Scan(&userRights) // Используйте Scan для извлечения значения

	// Проверка на наличие ошибки
	if err != nil {
		return "", err // Возвращаем пустую строку и ошибку
	}

	return userRights, nil // Возвращаем права пользователя и nil как ошибку
}

func (q *UserQueries) GetUserForResponsById(UserID uint) (models.UserResponse, error) {

	user := models.UserResponse{}

	query := "SELECT user_id, bd_used, first_name, last_name, patronymic_name, user_email, user_phone FROM users WHERE user_id = $1"

	err := q.Get(&user, query, UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("пользователь с ID %d не найден", UserID)
		}
		return user, err
	}

	rightsQuery := "SELECT user_rights FROM rights WHERE user_id = $1"

	var userRights string = ""

	err = q.Get(&userRights, rightsQuery, UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Оставляем значение по умолчанию, если права не найдены
		} else {
			return user, err
		}
	} else {
		user.UserRights = userRights
	}

	return user, nil
}

// Установка значения поля user_bd -
// бд с которой сейчас работает пользователь
func (q *UserQueries) SetupUserBd(id uint, user_bd uint) error {
	// Проверяем существование значения bd
	var _bd int

	query := "SELECT bd_used FROM users WHERE user_id = $1"
	err := q.Get(&_bd, query, id)

	if err != nil {
		return fmt.Errorf("ошибка при проверке существования bd: %w", err)
	}

	// Если значение bd уже существует и больше 0 и равно id, выходим из функции
	if _bd > 0 && _bd == int(id) {
		return nil
	}

	updateQuery := "UPDATE users SET bd_used = $1 WHERE user_id = $2"
	_, err = q.Exec(updateQuery, user_bd, id)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении bd_used: %w", err)
	}

	return nil
}
