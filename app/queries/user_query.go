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

	// Строка запроса, без поля UserID и CreatedUser.
	query := "INSERT INTO users ( bd_used, first_name, patronymic_name, last_name, user_email, user_phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := q.Exec(query, b.BdUsed, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.UserPhone, b.Password)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) SetupUserRight(id uint, userBd uint, role string) error {
	// Проверяем, что роль является допустимой
	// validRoles := map[string]bool{
	// 	"admin":    true,
	// 	"member":   true,
	// 	"director": true,
	// 	"manager":  true,
	// }
	// if !validRoles[role] {
	// 	return fmt.Errorf("недопустимая роль: %s", role)
	// }

	// Проверяем, существует ли уже такая запись
	existsQuery := `
		 SELECT 
			  COUNT(*) AS rights_count
		 FROM rights
		 WHERE user_id = $1 AND user_bd = $2 AND user_role = $3
	`

	var rightsCount int
	err := q.QueryRow(existsQuery, id, userBd, role).Scan(&rightsCount)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования записи: %w", err)
	}

	// Если права уже существуют, возвращаем ошибку
	if rightsCount > 0 {
		return fmt.Errorf("пользователь с ID %d уже имеет права %s для базы данных %d", id, role, userBd)
	}

	// Проверяем, является ли пользователь администратором в любой базе данных
	adminCheckQuery := `
		 SELECT COUNT(*) 
		 FROM rights 
		 WHERE user_id = $1 AND user_role = 'admin'
	`

	var adminCount int
	err = q.QueryRow(adminCheckQuery, id).Scan(&adminCount)
	if err != nil {
		return fmt.Errorf("ошибка проверки прав администратора: %w", err)
	}

	// Если пользователь уже является администратором в любой базе данных
	if adminCount > 0 {
		return fmt.Errorf("пользователь с ID %d уже является администратором в другой базе данных", id)
	}

	// Если пользователь пытается стать администратором, проверяем, что он не имеет прав в других базах данных
	if role == "admin" {
		otherRightsCheckQuery := `
			  SELECT COUNT(*) 
			  FROM rights 
			  WHERE user_id = $1 AND user_bd != $2
		 `

		var otherRightsCount int
		err = q.QueryRow(otherRightsCheckQuery, id, userBd).Scan(&otherRightsCount)
		if err != nil {
			return fmt.Errorf("ошибка проверки прав в других базах данных: %w", err)
		}

		// Если пользователь имеет права в других базах данных, он не может стать администратором
		if otherRightsCount > 0 {
			return fmt.Errorf("пользователь с ID %d имеет права в других базах данных и не может стать администратором", id)
		}
	}

	// Добавляем новые права
	query := "INSERT INTO rights (user_id, user_bd, holding, user_role) VALUES ($1, $2, $3, $4)"
	_, err = q.Exec(query, id, userBd, "", role) // holding передаётся как пустая строка
	if err != nil {
		return fmt.Errorf("ошибка при установке прав пользователя: %w", err)
	}

	return nil
}

// func (q *UserQueries) GetUserPasswordByEmail(UserEmail string) (string, error) {
// 	var passwordHash string

// 	query := "SELECT password FROM users WHERE user_email = $1"
// 	err := q.Get(&passwordHash, query, UserEmail)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", fmt.Errorf("пользователь с email %s не найден", UserEmail)
// 		}
// 		return "", err
// 	}

//		return passwordHash, nil
//	}
func (q *UserQueries) GetUserByEmail(UserEmail string) (models.UserType, error) {

	var user models.UserType

	// Запрос для получения данных пользователя
	query := "SELECT * FROM users WHERE user_email = $1"

	err := q.Get(&user, query, UserEmail)
	if err != nil {
		// Если пользователь не найден, возвращаем пустую структуру с email
		if err == sql.ErrNoRows {
			return models.UserType{UserEmail: UserEmail}, nil
		}
		// Если произошла другая ошибка, возвращаем её
		return models.UserType{}, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return user, nil
}

func (q *UserQueries) GetUserRightsById(UserID uint) ([]models.UserRights, error) {
	// Определяем переменную для хранения массива прав пользователя.
	var rights []models.UserRights

	// Определяем строку запроса.
	query := `
		 SELECT 
			  rights_id, 
			  created_rights, 
			  user_bd, 
			  holding, 
			  user_role 
		 FROM rights 
		 WHERE user_id = $1
	`

	// Выполняем запрос и сканируем результат в переменную rights.
	err := q.Select(&rights, query, UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если права не найдены, возвращаем пустой массив и nil (или ошибку, если это критично).
			return []models.UserRights{}, nil
		}
		// Возвращаем ошибку, если что-то пошло не так.
		return nil, fmt.Errorf("ошибка при получении прав пользователя: %w", err)
	}

	// Возвращаем массив прав пользователя.
	return rights, nil
}

// func (q *UserQueries) GetUserRightByID(userID uint) (string, error) {
// 	// Определяем переменную для хранения прав пользователя.
// 	var userRights string

// 	// Определяем строку запроса.
// 	query := "SELECT user_rights FROM rights WHERE user_id = $1"

// 	// Используйте QueryRow для получения строки.
// 	err := q.QueryRow(query, userID).Scan(&userRights) // Используйте Scan для извлечения значения

// 	// Проверка на наличие ошибки
// 	if err != nil {
// 		return "", err // Возвращаем пустую строку и ошибку
// 	}

// 	return userRights, nil // Возвращаем права пользователя и nil как ошибку
// }

// func (q *UserQueries) GetUserForResponsById(UserID uint) (models.UserType, error) {
// 	user := models.UserType{}

// 	// Запрос для получения данных пользователя без ненужных полей
// 	userQuery := `
// 		 SELECT
// 			  user_id,
// 			  bd_used,
// 			  first_name,
// 			  patronymic_name,
// 			  last_name,
// 			  user_email,
// 			  user_phone
// 		 FROM users
// 		 WHERE user_id = $1
// 	`

// 	err := q.Get(&user, userQuery, UserID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return user, fmt.Errorf("пользователь с ID %d не найден", UserID)
// 		}
// 		return user, err
// 	}

// 	// Запрос для получения прав пользователя
// 	rightsQuery := `
// 		 SELECT
// 			  user_rights,
// 			  user_bd,
// 			  holding
// 		 FROM rights
// 		 WHERE user_id = $1
// 	`

// 	var userRights struct {
// 		UserRights string `db:"user_rights"`
// 		UserBD     string `db:"user_bd"`
// 		Holding    string `db:"holding"`
// 	}

// 	err = q.Get(&userRights, rightsQuery, UserID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// Оставляем значения по умолчанию, если права не найдены
// 		} else {
// 			return user, err
// 		}
// 	} else {
// 		user.UserRights = userRights.UserRights
// 		user.UserBD = userRights.UserBD
// 		user.Holding = userRights.Holding
// 	}

// 	return user, nil
// }

// func (q *UserQueries) GetUserForResponsById(UserID uint) (models.UserType, error) {

// 	user := models.UserType{}

// 	query := "SELECT user_id, bd_used, first_name, patronymic_name, last_name, user_email, user_phone FROM users WHERE user_id = $1"

// 	err := q.Get(&user, query, UserID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return user, fmt.Errorf("пользователь с ID %d не найден", UserID)
// 		}
// 		return user, err
// 	}

// 	rightsQuery := "SELECT user_rights FROM rights WHERE user_id = $1"

// 	var userRights string = ""

// 	err = q.Get(&userRights, rightsQuery, UserID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// Оставляем значение по умолчанию, если права не найдены
// 		} else {
// 			return user, err
// 		}
// 	} else {
// 		user.UserRights = userRights
// 	}

// 	return user, nil
// }

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
