package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func CreateDBmain(db *sqlx.DB, db_new string) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname = $1);"

	// Проверка существования базы данных
	if err := db.Get(&exists, checkQuery, db_new); err != nil {
		log.Printf("Ошибка при проверке существования БД %s: %v", db_new, err)
		return err
	}

	if exists {
		log.Printf("База данных %s уже существует", db_new)
		return nil
	}

	// Создание базы данных
	createDBQuery := "CREATE DATABASE \"" + db_new + "\";"
	if _, err := db.Exec(createDBQuery); err != nil {
		log.Printf("Ошибка при создании БД %s: %v", db_new, err)
		return err
	}

	return nil
}
