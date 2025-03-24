package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Список пар "название таблицы" - "путь к SQL файлу".
var tablesStart = map[string]string{
	"users":  "platform/migrations/000001_create_users.up.sql",
	"rights": "platform/migrations/000002_create_rights.up.sql",
}

// CreateTables создает все таблицы из списка.
func CreateStartTables(db *sqlx.DB) {
	for tableName, sqlFile := range tablesStart {
		if err := executeSQLFile(db, sqlFile); err != nil {
			log.Printf("Ошибка при выполнении SQL файла для таблицы %s: %v", tableName, err)
			return
		}

		// Проверка успешного создания первой таблицы перед созданием второй
		if tableName == "users" {
			log.Printf("Таблица %s успешно создана. Переход к созданию следующей таблицы.", tableName)
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
