package queries

import (
	"ApiAyy/app/models"

	"github.com/jmoiron/sqlx"
)

type CompanyQueries struct {
	*sqlx.DB
}

// GetCompanies method for getting all companies.
// func (q *CompanyQueries) GetCompanies() ([]models.Company, error) {
// 	// Define Companies variable.
// 	Companies := []models.Company{}

// 	// Define query string.
// 	query := `SELECT * FROM Companies`

// 	// Send query to database.
// 	err := q.Select(&companies, query)
// 	if err != nil {
// 		// Return empty object and error.
// 		return companies, err
// 	}

// 	// Return query result.
// 	return companies, nil
// }

// GetCompany method for getting one company by given ID.
func (q *CompanyQueries) GetCompany(company_id uint) (models.Company, error) {
	// Define company variable.
	company := models.Company{}

	// Define query string.
	query := `SELECT * FROM company WHERE id_company = $1`

	// Send query to database.
	err := q.Get(&company, query, company_id)
	if err != nil {
		// Return empty object and error.
		return company, err
	}

	// Return query result.
	return company, nil
}

// CreateCompany method for creating Company by given Company object.
func (q *CompanyQueries) CreateCompany(b *models.Company) error {

	// Определите строку запроса, исключив поле USER_ID.
	query := "INSERT INTO company (created_at, status_company, company_name, inn, kpp, ogrn, ur_adress, mail_adress, phone, email, director) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, b.CreatedAt, b.Company_status, b.Company_name, b.Inn, b.Kpp, b.Ogrn, b.Ur_adress, b.Mail_adress, b.Phone, b.Email, b.Director)
	if err != nil {
		// Return only error.
		return err
	}

	return nil
}

// UpdateCompany method for updating company by given Company object.
func (q *CompanyQueries) UpdateCompany(company_id, b *models.Company) error {
	// Define query string.
	query := `UPDATE company SET company_status = $2, company_name = $3, inn = $4, kpp = $5, ogrn = $6, ur_adress = $7, mail_adress = $8, phone = $9, email = $10, director = $11 WHERE company_id = $1`

	// Send query to database.
	_, err := q.Exec(query, b.Company_Id, b.Company_status, b.Company_name, b.Inn, b.Kpp, b.Ogrn, b.Ur_adress, b.Mail_adress, b.Phone, b.Email, b.Director)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
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
