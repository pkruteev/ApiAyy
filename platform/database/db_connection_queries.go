package database

import "ApiAyy/app/queries"

// Queries struct for collect all app queries.
type Queries struct {
	// *queries.BookQueries // load queries from Book model
	*queries.UserQueries
	*queries.CompanyQueries
	*queries.MyUsersQueries
	// *queries.DataBaseQueries
}

// OpenDBConnection func for opening database connection.
func DBConnectionQueries(db_name string) (*Queries, error) {

	// Define a new PostgreSQL connection.
	db_main, err := DBConnection(db_name)
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:    &queries.UserQueries{DB: db_main},
		CompanyQueries: &queries.CompanyQueries{DB: db_main},
		MyUsersQueries: &queries.MyUsersQueries{DB: db_main},
		// DataBaseQueries: &queries.DataBaseQueries{DB: db_main},
	}, nil

}
