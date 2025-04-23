package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) RegisterUser(b *models.UserType) error {

	// Строка запроса, без поля UserID и CreatedUser.
	query := "INSERT INTO users ( bd_used, first_name, patronymic_name, last_name, user_email, email_verification, user_phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := q.Exec(query, b.BdUsed, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.EmailVerification, b.UserPhone, b.Password)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) SetupUserRight(id uint, userBd string, role string) error {
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
		 SELECT COUNT(*) AS rights_count
		 FROM rights
		 WHERE user_id = $1 AND user_bd = $2 AND user_role = $3
	`

	var rightsCount int
	err := q.QueryRow(existsQuery, id, userBd, role).Scan(&rightsCount)
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}

	// Если права уже существуют, возвращаем ошибку
	if rightsCount > 0 {
		return fmt.Errorf("пользователь с ID %d уже имеет права %s", id, role)
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

	// Если пользователь пытается стать администратором, проверяем, что его гнт в таблице rights
	if role == "admin" {
		// Проверяем наличие любых прав у пользователя
		rightsCheckQuery := `
			 SELECT COUNT(*) 
			 FROM rights 
			 WHERE user_id = $1
		`

		var rightsCount int
		err = q.QueryRow(rightsCheckQuery, id).Scan(&rightsCount)
		if err != nil {
			return fmt.Errorf("ошибка проверки прав пользователя: %w", err)
		}

		if rightsCount > 0 {
			return fmt.Errorf("пользователь с ID %d уже имеет права в системе", id)
		}
	}
	// Если все проверки пройдены, добавляем права
	insertQuery := `
		 INSERT INTO rights (user_id, user_bd, holding, user_role)
		 VALUES ($1, $2, $3, $4)
	`
	_, err = q.Exec(insertQuery, id, userBd, "", role)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении прав: %w", err)
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

// Функция для получения прав нескольких пользователей
func (q *UserQueries) GetUserRightsById(userID uint) ([]models.UserRights, error) {
	log.Printf("Начало выполнения GetUserRightsById для userID: %d", userID)
	var rights []models.UserRights

	query := `
		 SELECT 
			  rights_id,
			  user_id,
			  created_rights, 
			  user_bd, 
			  holding, 
			  user_role 
		 FROM rights 
		 WHERE user_id = $1`

	// Используем Select для получения нескольких записей
	err := q.Select(&rights, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Права не найдены для userID: %d", userID)
			return []models.UserRights{}, nil
		}
		log.Printf("Ошибка при выполнении запроса для userID %d: %v", userID, err)
		return nil, fmt.Errorf("ошибка при получении прав пользователя: %w", err)
	}

	log.Printf("Успешно получено %d прав для userID: %d", len(rights), userID)
	if len(rights) > 0 {
		log.Printf("Первая запись прав: %+v", rights[0])
	}

	return rights, nil
}

// Функция получения прав для одного пользователя
func (q *UserQueries) GetUserRights(userID uint) (*models.UserRights, error) {
	log.Printf("Получение прав для userID: %d", userID)
	var rights models.UserRights

	query := `
		 SELECT 
			  rights_id,
			  user_id,
			  created_rights,
			  user_bd,
			  holding,
			  user_role
		 FROM rights 
		 WHERE user_id = $1
		 LIMIT 1`

	err := q.Get(&rights, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Права не найдены для userID: %d", userID)
			return nil, nil
		}
		log.Printf("Ошибка запроса прав: %v", err)
		return nil, fmt.Errorf("ошибка получения прав: %w", err)
	}

	log.Printf("Получены права: %+v", rights)
	return &rights, nil
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
func (q *UserQueries) SetupUserBd(id uint, user_bd string) error {
	var currentBd string

	query := "SELECT bd_used FROM users WHERE user_id = $1"
	err := q.Get(&currentBd, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("пользователь с id %d не найден", id)
		}
		return fmt.Errorf("ошибка при проверке bd_used: %w", err)
	}

	if currentBd == user_bd {
		return nil
	}

	updateQuery := "UPDATE users SET bd_used = $1 WHERE user_id = $2"
	_, err = q.Exec(updateQuery, user_bd, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении bd_used: %w", err)
	}

	return nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (q *UserQueries) IsAdmin(userID uint) (bool, error) {
	adminCheckQuery := `
		 SELECT COUNT(*) > 0
		 FROM rights 
		 WHERE user_id = $1 AND user_role = 'admin'
	`

	var isAdmin bool
	err := q.QueryRow(adminCheckQuery, userID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки прав администратора: %w", err)
	}

	return isAdmin, nil
}

// GetUserHolding проверяет наличие и возвращает значение holding для пользователя
func (q *UserQueries) GetUserHolding(userID uint) (string, error) {
	query := `
		 SELECT holding 
		 FROM rights 
		 WHERE user_id = $1
	`

	var holding string
	err := q.QueryRow(query, userID).Scan(&holding)

	if err == sql.ErrNoRows {
		return "", nil // Нет записи - возвращаем пустую строку
	}
	if err != nil {
		return "", fmt.Errorf("ошибка при получении holding: %w", err)
	}

	return holding, nil
}

// UpsertUserHolding обновляет или создает запись holding для пользователя
func (q *UserQueries) UpsertUserHolding(userID uint, holding string) error {
	query := `
		 INSERT INTO rights (user_id, holding)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id) 
		 DO UPDATE SET holding = EXCLUDED.holding
	`

	_, err := q.Exec(query, userID, holding)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении holding: %w", err)
	}

	return nil
}

func (q *UserQueries) UpdateUserHolding(userID uint, holding string) error {
	const query = `
		 UPDATE rights 
		 SET holding = $1 
		 WHERE user_id = $2`

	_, err := q.Exec(query, holding, userID)
	return err
}
