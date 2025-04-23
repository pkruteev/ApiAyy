package queries

import (
	"ApiAyy/app/models"
	"database/sql"

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
		created_ob, 
		typereal, 
		city, 
		street, 
		house, 
		flat, 
		floor, 
		square
	) VALUES (
		CURRENT_TIMESTAMP, 
		$1, $2, $3, $4, $5, $6, $7
	) RETURNING object_id`

	// Используем QueryRow для получения ID созданного объекта
	err := q.QueryRow(
		query,
		b.TypeReal,
		b.City,
		b.Street,
		b.House,
		b.Flat,
		b.Floor,
		b.Square,
	).Scan(&b.ObjectId)

	if err != nil {
		return err
	}

	return nil
}

// UpdateObject обновляет данные объекта недвижимости
func (q *ObjectQueries) UpdateObject(objectID uint, b *models.Objects) error {
	query := `UPDATE objects SET 
		typereal = $1, 
		city = $2, 
		street = $3, 
		house = $4, 
		flat = $5, 
		floor = $6, 
		square = $7 
	WHERE object_id = $8`

	_, err := q.Exec(
		query,
		b.TypeReal,
		b.City,
		b.Street,
		b.House,
		b.Flat,
		b.Floor,
		b.Square,
		objectID,
	)

	return err
}

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
