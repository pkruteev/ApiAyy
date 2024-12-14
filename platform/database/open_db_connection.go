package database

import "ApiAyy/app/queries"

// Queries struct for collect all app queries.
type Queries struct {
	// *queries.BookQueries // load queries from Book model
	*queries.UserQueries
	*queries.CompanyQueries
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		UserQueries:    &queries.UserQueries{DB: db},    // from User model
		CompanyQueries: &queries.CompanyQueries{DB: db}, // from Company model
	}, nil

	// return &Queries{
	// 	// Set queries from models:
	// 	BookQueries: &queries.BookQueries{DB: db}, // from Book model
	// }, nil
}
