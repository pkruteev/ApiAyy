package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func CreateDB(db_name string, db *sqlx.DB) error {

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname = $1);"
	if err := db.Get(&exists, checkQuery, db_name); err != nil {
		log.Printf("Ошибка при проверке существования БД %s: %v", db_name, err)
		return err
	}

	if exists {
		log.Printf("База данных %s уже существует", db_name)
		return nil
	}

	createDBQuery := "CREATE DATABASE \"" + db_name + "\";"
	if _, err := db.Exec(createDBQuery); err != nil {
		log.Printf("Ошибка при создании БД %s: %v", db_name, err)
		return err
	}
	return nil
}
