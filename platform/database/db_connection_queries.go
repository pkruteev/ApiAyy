package database

import (
	"ApiAyy/app/queries"

	"github.com/jmoiron/sqlx"
)

// Queries struct for collect all app queries.
type Queries struct {
	*queries.ObjectQueries
	*queries.UserQueries
	*queries.CompanyQueries
	*queries.MyUsersQueries
	DB *sqlx.DB
}

// Close closes the database connection
func (q *Queries) Close() error {
	if q.DB != nil {
		return q.DB.Close()
	}
	return nil
}

// OpenDBConnection func for opening database connection.
func DBConnectionQueries(db_name string) (*Queries, error) {

	// Define a new PostgreSQL connection.
	db_main, err := DBConnection(db_name)
	if err != nil {
		return nil, err
	}

	return &Queries{
		ObjectQueries:  &queries.ObjectQueries{DB: db_main},
		UserQueries:    &queries.UserQueries{DB: db_main},
		CompanyQueries: &queries.CompanyQueries{DB: db_main},
		MyUsersQueries: &queries.MyUsersQueries{DB: db_main},
	}, nil
}
