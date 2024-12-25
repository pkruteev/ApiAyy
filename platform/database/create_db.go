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

// var bd = map[string]string{
// 	"new_table": "platform/migrations/000003_create_company_table.up.sql",
// }

// CreateDatabase создает БД по userId и возвращает ошибку, если она возникла.
// func CreateDB(userId string, db_main *sqlx.DB) error {
// 	// Создаем SQL запрос для создания базы данных
// 	createDBQuery := "CREATE DATABASE \"" + userId + "\";"
// 	//Выполняем создание базы данных
// 	if _, err := db_main.Exec(createDBQuery); err != nil {
// 		log.Printf("Ошибка при создании БД %s: %v", userId, err)
// 		return err
// 	}

// 	return nil
// }

// takeSQLFile выполняет SQL-запрос из указанного файла.
// func takeSQLFile(db_new *sqlx.DB, filePath string) error {
// 	sqlQuery, err := os.ReadFile(filePath)
// 	if err != nil {
// 		log.Printf("Не удалось прочитать SQL файл %s: %v", filePath, err)
// 		return err
// 	}

// 	_, err = db_new.Exec(string(sqlQuery))
// 	if err != nil {
// 		log.Printf("Ошибка выполнения SQL: %v", err)
// 	}
// 	return err
// }
