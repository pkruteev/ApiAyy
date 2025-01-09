package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func CreateDB(userId string, db_main *sqlx.DB) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname = $1);"
	if err := db_main.Get(&exists, checkQuery, userId); err != nil {
		log.Printf("Ошибка при проверке существования БД %s: %v", userId, err)
		return err
	}

	if exists {
		log.Printf("База данных %s уже существует", userId)
		return nil
	}

	createDBQuery := "CREATE DATABASE \"" + userId + "\";"
	if _, err := db_main.Exec(createDBQuery); err != nil {
		log.Printf("Ошибка при создании БД %s: %v", userId, err)
		return err
	}

	return nil
}
