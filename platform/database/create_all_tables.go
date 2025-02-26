package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Список пар "название таблицы" - "путь к SQL файлу".
var all_tables = map[string]string{
	"companies":              "platform/migrations/000003_create_companies_table.up.sql",
	"r_schets":               "platform/migrations/000004_create_r_schets_table.up.sql",
	"contragencies":          "platform/migrations/000006_create_contragencies_table.up.sql",
	"r_schets_contragencies": "platform/migrations/000007_create_contragencies_r_schets_table.up.sql",
	"objects":                "platform/migrations/000005_create_objects_table.up.sql",
	"contracts":              "platform/migrations/000008_create_contracts_table.up.sql",
}

// CreateAllTables создает все таблицы из списка.
func CreateAllTables(db *sqlx.DB) {
	tableNames := make([]string, 0, len(all_tables))
	for tableName := range all_tables {
		tableNames = append(tableNames, tableName)
	}

	for i := 0; i < len(tableNames); i++ {
		tableName := tableNames[i]
		sqlFile := all_tables[tableName]

		log.Printf("Начинаем создание таблицы %s из файла %s", tableName, sqlFile)

		if err := executeAllSQLFile(db, sqlFile); err != nil {
			log.Printf("Ошибка при выполнении SQL файла для таблицы %s: %v", tableName, err)
			return
		}

		log.Printf("Таблица %s успешно создана.", tableName)
	}
}

// executeAllSQLFile выполняет SQL-запрос из указанного файла.
func executeAllSQLFile(db *sqlx.DB, filePath string) error {
	sqlQuery, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Не удалось прочитать SQL файл %s: %v", filePath, err)
		return err
	}

	log.Printf("Выполнение SQL для файла %s", filePath)
	_, err = db.Exec(string(sqlQuery))
	if err != nil {
		log.Printf("Ошибка выполнения SQL для файла %s: %v", filePath, err)
	}
	return err
}
