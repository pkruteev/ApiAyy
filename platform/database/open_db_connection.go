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
	db_main, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		UserQueries:    &queries.UserQueries{DB: db_main},    // from User model
		CompanyQueries: &queries.CompanyQueries{DB: db_main}, // from Company model
	}, nil

	// return &Queries{
	// 	// Set queries from models:
	// 	BookQueries: &queries.BookQueries{DB: db}, // from Book model
	// }, nil
}
