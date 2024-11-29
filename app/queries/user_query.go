package queries

import (
	"ApiAyy/app/models"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

// GetUsers method for getting all users.
// func (q *UserQueries) GetUsers() ([]models.User, error) {
// 	// Define books variable.
// 	users := []models.User{}

// 	// Define query string.
// 	query := `SELECT * FROM users`

// 	// Send query to database.
// 	err := q.Select(&users, query)
// 	if err != nil {
// 		// Return empty object and error.
// 		return users, err
// 	}

// 	// Return query result.
// 	return users, nil
// }

// GetBook method for getting one book by given ID.
// func (q *UserQueries) GetUser(id uuid.UUID) (models.Book, error) {
// 	// Define book variable.
// 	book := models.Book{}

// 	// Define query string.
// 	query := `SELECT * FROM books WHERE id = $1`

// 	// Send query to database.
// 	err := q.Get(&book, query, id)
// 	if err != nil {
// 		// Return empty object and error.
// 		return book, err
// 	}

// 	// Return query result.
// 	return book, nil
// }

// CreateBook method for creating book by given Book object.
func (q *UserQueries) RegisterUser(b *models.UserType) error {
	// Define query string.
	query := `INSERT INTO users VALUES ( $2, $3, $4, $5, $6, $7, $7, $9)`

	// Отправьте запрос в базу данных.
	_, err := q.Exec(query, b.FirstName, b.PatronymicName, b.LastName, b.UserEmail, b.UserPhone, b.UserCompany, b.Password, b.UserRole)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateBook method for updating book by given Book object.
// func (q *BookQueries) UpdateBook(id uuid.UUID, b *models.Book) error {
// 	// Define query string.
// 	query := `UPDATE books SET updated_at = $2, title = $3, author = $4, book_status = $5, book_attrs = $6 WHERE id = $1`

// 	// Send query to database.
// 	_, err := q.Exec(query, id, b.UpdatedAt, b.Title, b.Author, b.BookStatus, b.BookAttrs)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// DeleteBook method for delete book by given ID.
// func (q *BookQueries) DeleteBook(id uuid.UUID) error {
// 	// Define query string.
// 	query := `DELETE FROM books WHERE id = $1`

// 	// Send query to database.
// 	_, err := q.Exec(query, id)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }
