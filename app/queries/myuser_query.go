package queries

import (
	"ApiAyy/app/models"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// MyUsersQueries struct for queries from User model.
type MyUsersQueries struct {
	DB *sqlx.DB
}

func (q *MyUsersQueries) GetMyAllUsers(userBd string) ([]models.MyUserType, error) {
	type rightsRecord struct {
		RightsID uint   `db:"rights_id"`
		UserID   uint   `db:"user_id"`
		UserRole string `db:"user_role"`
	}

	log.Printf("Начало GetMyAllUsers для userBd: %s", userBd)

	// 1. Получаем данные из таблицы rights
	var rightsRecords []rightsRecord
	rightsQuery := "SELECT rights_id, user_id, user_role FROM rights WHERE user_bd = $1"
	err := q.DB.Select(&rightsRecords, rightsQuery, userBd)
	if err != nil {
		log.Printf("Ошибка при запросе к таблице rights: %v", err)
		return nil, fmt.Errorf("ошибка при получении user_id из таблицы rights: %w", err)
	}

	if len(rightsRecords) == 0 {
		return []models.MyUserType{}, nil
	}

	// 2. Получаем ID пользователей
	userIds := make([]uint, len(rightsRecords))
	for i, record := range rightsRecords {
		userIds[i] = record.UserID
	}

	// 3. Получаем данные пользователей
	var users []models.MyUserType
	usersQuery := `
		 SELECT 
			  user_id, 
			  first_name, 
			  patronymic_name, 
			  last_name, 
			  user_email, 
			  user_phone
		 FROM users 
		 WHERE user_id = ANY($1)
	`
	err = q.DB.Select(&users, usersQuery, pq.Array(userIds))
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных пользователей: %w", err)
	}

	// 4. Объединяем данные
	result := make([]models.MyUserType, 0, len(rightsRecords))
	userMap := make(map[uint]models.MyUserType, len(users))

	// Сначала заполняем мапу пользователями
	for _, user := range users {
		userMap[user.UserID] = user
	}

	// Затем добавляем данные из прав
	for _, rightsRec := range rightsRecords {
		if user, exists := userMap[rightsRec.UserID]; exists {
			// Копируем все нужные поля из прав
			user.RightsID = rightsRec.RightsID // Добавляем RightsID
			user.UserRole = rightsRec.UserRole
			result = append(result, user)
		}
	}

	return result, nil
}

// проверить существование записи по RightsID
func (q *UserQueries) RightsExists(rightsID uint) (bool, error) {
	log.Printf("Проверка существования прав с rightsID: %d", rightsID)

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rights WHERE rights_id = $1)`

	err := q.DB.Get(&exists, query, rightsID)
	if err != nil {
		log.Printf("Ошибка при проверке существования прав %d: %v", rightsID, err)
		return false, fmt.Errorf("ошибка при проверке существования прав: %w", err)
	}

	if !exists {
		log.Printf("Права с rightsID %d не найдены", rightsID)
		return false, nil
	}

	log.Printf("Права с rightsID %d существуют", rightsID)
	return true, nil
}

func (q *UserQueries) DeleteRightsMyUser(RightsId uint) error {

	query := `DELETE FROM rights WHERE rights_id = $1`

	_, err := q.Exec(query, RightsId)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
