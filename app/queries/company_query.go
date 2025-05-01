package queries

import (
	"ApiAyy/app/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CompanyQueries struct {
	*sqlx.DB
}

// GetCompanies возвращает список всех компаний (counterparty = false)
func (q *CompanyQueries) GetCompanies() ([]models.Company, error) {
	companies := []models.Company{}
	query := `SELECT * FROM companies WHERE is_counterparty = false ORDER BY company_id`
	err := q.Select(&companies, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get companies: %w", err)
	}
	return companies, nil
}

// GetCompany возвращает одну компанию по ID
func (q *CompanyQueries) GetCompany(id uint) (models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE company_id = $1 AND counterparty = false LIMIT 1`
	err := q.Get(&company, query, id)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to get company: %w", err)
	}
	return company, nil
}

// CreateCompany создает новую компанию (counterparty = false)
func (q *CompanyQueries) CreateCompany(c *models.Company) error {
	c.IsCounterparty = false // Явно устанавливаем флаг компании

	query := `INSERT INTO companies (
		is_counterparty, status, name, inn, kpp, ogrn, data_ogrn, 
		ogrnip, data_ogrnip, ur_address, mail_address, phone, email, director
	) VALUES (
		:is_counterparty, :status, :name, :inn, :kpp, :ogrn, :data_ogrn, 
		:ogrnip, :data_ogrnip, :ur_address, :mail_address, :phone, :email, :director
	) RETURNING company_id`

	_, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// UpdateCompany обновляет данные компании
func (q *CompanyQueries) UpdateCompany(id uint, c *models.Company) error {
	c.IsCounterparty = false // Гарантируем что это компания
	c.CompanyId = id         // Устанавливаем ID из параметра

	query := `UPDATE companies SET 
		status = :status, name = :name, kpp = :kpp, ogrn = :ogrn, 
		data_ogrn = :data_ogrn, ogrnip = :ogrnip, data_ogrnip = :data_ogrnip, 
		ur_address = :ur_address, mail_address = :mail_address, 
		phone = :phone, email = :email, director = :director
	WHERE company_id = :company_id AND is_counterparty = false`

	result, err := q.NamedExec(query, c)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
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

// CompanyExists проверяет существование компании
func (q *CompanyQueries) CompanyExists(id uint) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE company_id = $1 AND is_counterparty = false)`
	err := q.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check company existence: %w", err)
	}
	return exists, nil
}

// DeleteCompany удаляет компанию по ID
func (q *CompanyQueries) DeleteCompany(id uint) error {
	query := `DELETE FROM companies WHERE company_id = $1 AND is_counterparty = false`
	result, err := q.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
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
