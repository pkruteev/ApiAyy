package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Список пар "название таблицы" - "путь к SQL файлу".
var tables = map[string]string{

	"users":   "platform/migrations/000001_create_users_table.up.sql",
	"company": "platform/migrations/000002_create_company_table.up.sql",

	// "company": "platform/migrations/000003_create_companies_table.up.sql",
}

// CreateTables создает все таблицы из списка.
func CreateTables(db *sqlx.DB) {
	for tableName, sqlFile := range tables {
		if err := executeSQLFile(db, sqlFile); err != nil {
			log.Printf("Ошибка при выполнении SQL файла для таблицы %s: %v", tableName, err)
		}
	}
}

// executeSQLFile выполняет SQL-запрос из указанного файла.
func executeSQLFile(db *sqlx.DB, filePath string) error {
	sqlQuery, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Не удалось прочитать SQL файл %s: %v", filePath, err)
		return err
	}

	_, err = db.Exec(string(sqlQuery))
	if err != nil {
		log.Printf("Ошибка выполнения SQL: %v", err)
	}
	return err
}
