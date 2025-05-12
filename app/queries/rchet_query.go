package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RSchetQueries struct {
	*sqlx.DB
}

// GetRSchets возвращает список всех расчетных счетов
func (q *RSchetQueries) GetRSchets() ([]models.RSchet, error) {
	rschets := []models.RSchet{}
	query := `SELECT * FROM r_schets ORDER BY r_schet_id`
	err := q.Select(&rschets, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get r_schets: %w", err)
	}
	return rschets, nil
}

// GetRSchet возвращает один расчетный счет по ID
func (q *RSchetQueries) GetRSchet(id uint) (models.RSchet, error) {
	var rschet models.RSchet
	query := `SELECT * FROM r_schets WHERE r_schet_id = $1 LIMIT 1`
	err := q.Get(&rschet, query, id)
	if err != nil {
		return models.RSchet{}, fmt.Errorf("failed to get r_schet: %w", err)
	}
	return rschet, nil
}

// CreateRSchet создает новый расчетный счет
func (q *RSchetQueries) CreateRSchet(r *models.RSchet) error {
	query := `INSERT INTO r_schets (
		company_id, created_schet, r_schet, bank_name, bank_bic, kor_schet
	) VALUES (
		:company_id, :created_schet, :r_schet, :bank_name, :bank_bic, :kor_schet
	) RETURNING r_schet_id`

	_, err := q.NamedExec(query, r)
	if err != nil {
		return fmt.Errorf("failed to create r_schet: %w", err)
	}
	return nil
}

// UpdateRSchet обновляет данные расчетного счета
func (q *RSchetQueries) UpdateRSchet(id uint, r *models.RSchet) error {
	r.RSchetId = id // Устанавливаем ID из параметра

	query := `UPDATE r_schets SET 
		company_id = :company_id,
		created_schet = :created_schet, 
		r_schet = :r_schet, 
		bank_name = :bank_name, 
		bank_bic = :bank_bic, 
		kor_schet = :kor_schet
	WHERE r_schet_id = :r_schet_id`

	result, err := q.NamedExec(query, r)
	if err != nil {
		return fmt.Errorf("failed to update r_schet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// RSchetExists проверяет существование расчетного счета
func (q *RSchetQueries) RSchetExists(id uint) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM r_schets WHERE r_schet_id = $1)`
	err := q.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check r_schet existence: %w", err)
	}
	return exists, nil
}

// DeleteRSchet удаляет расчетный счет по ID
func (q *RSchetQueries) DeleteRSchet(id uint) error {
	query := `DELETE FROM r_schets WHERE r_schet_id = $1`
	result, err := q.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete r_schet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteRSchets удаляет несколько расчетных счетов по IDs
func (q *RSchetQueries) DeleteRSchets(ids []uint) error {
	query := `DELETE FROM r_schets WHERE r_schet_id = ANY($1)`
	result, err := q.Exec(query, ids)
	if err != nil {
		return fmt.Errorf("failed to delete r_schets: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
