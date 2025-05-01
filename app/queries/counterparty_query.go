package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CounterpartyQueries struct {
	*sqlx.DB
}

// GetCounterparties возвращает список всех контрагентов (counterparty = true)
func (q *CounterpartyQueries) GetCounterparties() ([]models.Company, error) {
	counterparties := []models.Company{}
	query := `SELECT * FROM companies WHERE is_counterparty = true ORDER BY company_id`
	err := q.Select(&counterparties, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get counterparties: %w", err)
	}
	return counterparties, nil
}

// GetCounterparty возвращает одного контрагента по ID
func (q *CounterpartyQueries) GetCounterparty(id uint) (models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE company_id = $1 AND is_counterparty = true LIMIT 1`
	err := q.Get(&company, query, id)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to get counterparty: %w", err)
	}
	return company, nil
}

// CreateCounterparty создает нового контрагента
func (q *CounterpartyQueries) CreateCounterparty(c *models.Company) error {
	c.IsCounterparty = true // Устанавливаем флаг контрагента

	query := `INSERT INTO companies (
		is_counterparty, status, name, inn, kpp, ogrn, data_ogrn, 
		ogrnip, data_ogrnip, ur_address, mail_address, phone, email, director
	) VALUES (
		:is_counterparty, :status, :name, :inn, :kpp, :ogrn, :data_ogrn, 
		:ogrnip, :data_ogrnip, :ur_address, :mail_address, :phone, :email, :director
	) RETURNING company_id`

	_, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to create counterparty: %w", err)
	}
	return nil
}

// UpdateCounterparty обновляет данные контрагента
func (q *CounterpartyQueries) UpdateCounterparty(id uint, c *models.Company) error {
	c.IsCounterparty = true // Гарантируем что это контрагент
	c.CompanyId = id        // Устанавливаем ID из параметра

	query := `UPDATE companies SET 
		status = :status, name = :name, kpp = :kpp, ogrn = :ogrn, 
		data_ogrn = :data_ogrn, ogrnip = :ogrnip, data_ogrnip = :data_ogrnip, 
		ur_address = :ur_address, mail_address = :mail_address, 
		phone = :phone, email = :email, director = :director
	WHERE company_id = :company_id AND is_counterparty = true`

	result, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to update counterparty: %w", err)
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

// CounterpartyExists проверяет существование контрагента
func (q *CounterpartyQueries) CounterpartyExists(id uint) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE company_id = $1 AND is_counterparty = true)`
	err := q.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check counterparty existence: %w", err)
	}
	return exists, nil
}

// DeleteCounterparty удаляет контрагента по ID
func (q *CounterpartyQueries) DeleteCounterparty(id uint) error {
	query := `DELETE FROM companies WHERE company_id = $1 AND is_counterparty = true`
	result, err := q.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete counterparty: %w", err)
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
