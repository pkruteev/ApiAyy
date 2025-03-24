package queries

import (
	"ApiAyy/app/models"

	"github.com/jmoiron/sqlx"
)

type CompanyQueries struct {
	*sqlx.DB
}

// GetCompanies method for getting all companies.
func (q *CompanyQueries) GetCompanies() ([]models.Company, error) {
	// Определяем переменную для хранения компаний.
	companies := []models.Company{}

	// Строка запроса для таблицы companies.
	query := "SELECT * FROM companies"

	// Используем Select для получения всех строк.
	err := q.Select(&companies, query)
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}

	// Если запрос выполнен успешно, возвращаем список компаний (пустой или нет).
	return companies, nil
}

// GetCompany method for getting one company by given ID.
func (q *CompanyQueries) GetCompany() ([]models.Company, error) {
	// Define company variable.
	companies := []models.Company{}

	// Define query string.
	query := `SELECT * FROM companies`

	// Send query to database.
	err := q.Select(&companies, query)
	if err != nil {
		// Return empty object and error.
		return companies, err
	}

	// Return query result.
	return companies, nil
}

// CreateCompany method for creating Company by given Company object.
func (q *CompanyQueries) CreateCompany(b *models.Company) error {
	// Определите строку запроса, исключив поле Company_Id.
	query := "INSERT INTO companies (contragent, status, name, inn, kpp, ogrn, data_ogrn, ogrnip, data_ogrnip, ur_address, mail_address, phone, email, director) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, b.Contragent, b.Status, b.Name, b.Inn, b.Kpp, b.Ogrn, b.DataOgrn, b.Ogrnip, b.DataOgrnip, b.UrAddress, b.MailAddress, b.Phone, b.Email, b.Director)
	if err != nil {
		// Вернуть только ошибку.
		return err
	}

	return nil
}

// UpdateCompany method for updating company by given Company object.
func (q *CompanyQueries) UpdateCompany(company_id, b *models.Company) error {
	// Define query string.
	query := `UPDATE company SET contragent = $2, status = $3, name = $4, inn = $5, kpp = $6, ogrn = $7, data_ogrn = $8, ogrnip= $9, data_ogrnip= $10, ur_address = $11, mail_address = $12, phone = $13, email = $14, director = $15 WHERE company_id = $1`

	// Send query to database.
	_, err := q.Exec(query, b.Contragent, b.Status, b.Name, b.Inn, b.Kpp, b.Ogrn, b.DataOgrn, b.Ogrnip, b.DataOgrnip, b.UrAddress, b.MailAddress, b.Phone, b.Email, b.Director)
	if err != nil {
		// Return only error.
		return err
	}

	return nil
}

// DeleteBook method for delete book by given ID.
func (q *CompanyQueries) DeleteCompany(company_id int) error {
	// Define query string.
	query := `DELETE FROM company WHERE company_id = $1`

	// Send query to database.
	_, err := q.Exec(query, company_id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
