package database

import "ApiAyy/app/queries"

// Queries struct for collect all app queries.
type UserQueries struct {
	*queries.UserQueries
}

func DBMainConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db_main, err := DBConnection("main")
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db_main},
	}, nil

}
