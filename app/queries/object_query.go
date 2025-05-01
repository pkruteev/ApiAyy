package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ObjectQueries struct {
	*sqlx.DB
}

// GetObjects возвращает список всех объектов недвижимости
func (q *ObjectQueries) GetObjects() ([]models.Objects, error) {
	var objects []models.Objects
	query := `SELECT * FROM objects ORDER BY object_id`
	err := q.Select(&objects, query)
	return objects, err
}

// GetObject возвращает один объект по ID
func (q *ObjectQueries) GetObject(id uint) (*models.Objects, error) {
	object := &models.Objects{}
	query := `SELECT * FROM objects WHERE object_id = $1`
	err := q.Get(object, query, id)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// CreateObject создает новый объект недвижимости
func (q *ObjectQueries) CreateObject(b *models.Objects) error {
	query := `INSERT INTO objects (
		 cadastr_number,
		 typereal, 
		 city, 
		 street, 
		 house, 
		 flat, 
		 floor, 
		 square,
		 company_id
	) VALUES (
		 $1, $2, $3, $4, $5, $6, $7, $8, $9
	) RETURNING object_id`

	err := q.QueryRow(
		query,
		b.CadastrNumber,
		b.TypeReal,
		b.City,
		b.Street,
		b.House,
		b.Flat,
		b.Floor,
		b.Square,
		b.CompanyId,
	).Scan(&b.ObjectId)

	return err
}

// UpdateObject обновляет данные объекта недвижимости
func (q *ObjectQueries) UpdateObject(objectID uint, b *models.Objects) error {

	query := `UPDATE objects SET 
		 company_id = $1,
		 cadastr_number = $2,
		 typereal = $3, 
		 city = $4, 
		 street = $5, 
		 house = $6, 
		 flat = $7, 
		 floor = $8, 
		 square = $9
	WHERE object_id = $10`

	_, err := q.Exec(
		query,
		b.CompanyId,
		b.CadastrNumber,
		b.TypeReal,
		b.City,
		b.Street,
		b.House,
		b.Flat,
		b.Floor,
		b.Square,
		objectID,
	)

	if err != nil {
		return fmt.Errorf("ошибка обновления объекта: %w", err)
	}

	return nil
}

// CompanyExists проверяет существование компании в базе данных
// func (q *ObjectQueries) CompanyExists(companyID uint) (bool, error) {
// 	var exists bool
// 	err := q.QueryRow(
// 		"SELECT EXISTS(SELECT 1 FROM companies WHERE company_id = $1)",
// 		companyID,
// 	).Scan(&exists)

// 	if err != nil {
// 		return false, fmt.Errorf("ошибка проверки существования компании: %w", err)
// 	}

// 	return exists, nil
// }

// ObjectExists проверяет существование объекта по ID
func (q *ObjectQueries) ObjectExists(objectID uint) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM objects WHERE object_id = $1)`
	var exists bool
	err := q.Get(&exists, query, objectID)
	return exists, err
}

// DeleteObject удаляет объект по ID
func (q *ObjectQueries) DeleteObject(objectID uint) error {
	query := `DELETE FROM objects WHERE object_id = $1`
	result, err := q.Exec(query, objectID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
