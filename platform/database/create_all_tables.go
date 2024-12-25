package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Список пар "название таблицы" - "путь к SQL файлу".
var all_tables = map[string]string{

	"companies": "platform/migrations/000003_create_companies_table.up.sql",
	"r_schets":  "platform/migrations/000004_create_r_schets_table.up.sql",
	"objects":   "platform/migrations/000005_create_objects_table.up.sql",
}

// CreateTables создает все таблицы из списка.
func CreateAllTables(db_new *sqlx.DB) {
	for tableName, sqlFile := range all_tables {
		if err := executeAllSQLFile(db_new, sqlFile); err != nil {
			log.Printf("Ошибка при выполнении SQL файла для таблицы %s: %v", tableName, err)
		}
	}
}

// executeSQLFile выполняет SQL-запрос из указанного файла.
func executeAllSQLFile(db_new *sqlx.DB, filePath string) error {
	sqlQuery, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Не удалось прочитать SQL файл %s: %v", filePath, err)
		return err
	}

	_, err = db_new.Exec(string(sqlQuery))
	if err != nil {
		log.Printf("Ошибка выполнения SQL: %v", err)
	}
	return err
}
