package queries

import (
	"ApiAyy/app/models"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type MyUsersQueries struct {
	*sqlx.DB
}

func (q *MyUsersQueries) GetMyUsers(userId uint) ([]models.UserType, error) {
	// Определяем переменную для хранения пользователей.
	users := []models.UserType{}

	// Строка запроса для таблицы rights.
	query := "SELECT * FROM rights WHERE user_bd = $1"

	// Выводим userId для отладки.
	// log.Printf("Executing query with userId: %d", userId)

	// Используйте Select для получения всех строк.
	err := q.Select(&users, query, userId)
	if err != nil {
		return users, err
	}
	// Выводим найденных пользователей для отладки.
	// log.Printf("Retrieved users: %+v", users)
	return users, nil
}

// func (q *BookQueries) GetBooks() ([]models.Book, error) {
// 	// Define books variable.
// 	books := []models.Book{}

// 	// Define query string.
// 	query := `SELECT * FROM books`

// 	// Send query to database.
// 	err := q.Select(&books, query)
// 	if err != nil {
// 		// Return empty object and error.
// 		return books, err
// 	}

// 	// Return query result.
// 	return books, nil
// }

// GetBook method for getting one book by given ID.
// func (q *BookQueries) GetBook(id uuid.UUID) (models.Book, error) {
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
// func (q *BookQueries) CreateBook(b *models.Book) error {
// 	// Define query string.
// 	query := `INSERT INTO books VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

// 	// Send query to database.
// 	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Author, b.BookStatus, b.BookAttrs)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }
