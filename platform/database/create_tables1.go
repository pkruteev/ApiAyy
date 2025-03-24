package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// loadSQL загружает SQL-запрос из файла
func loadSQL1(filename string) (string, error) {
	// Указываем абсолютный путь к папке migrations
	basePath := `C:\IT\ApiAyy\platform\migrations`

	// Формируем полный путь к файлу
	filePath := filepath.Join(basePath, filename)

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("файл не найден: " + filePath)
	}

	// Читаем содержимое файла
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Возвращаем содержимое файла в виде строки
	return string(content), nil
}

// CreateTables создает все таблицы в базе данных
func CreateTables1(db *sqlx.DB) error {
	// Установка временной зоны
	if _, err := db.Exec("SET TIMEZONE='Europe/Moscow';"); err != nil {
		return err
	}

	// Список SQL-файлов для создания таблиц
	sqlFiles := []string{
		"000003_create_companies_table.up.sql",
		"000004_create_r_schets_table.up.sql",
		"000005_create_objects_table.up.sql",
		"000006_create_contragencies_table.up.sql",
		"000007_create_contragencies_r_schets_table.up.sql",
		"000008_create_contracts_table.up.sql",
	}

	// Выполнение SQL-запросов в строгом порядке
	for _, file := range sqlFiles {
		query, err := loadSQL(file)
		if err != nil {
			// Логируем ошибку, но продолжаем выполнение
			log.Printf("Ошибка загрузки SQL-запроса из файла %s: %v", file, err)
			return err // Возвращаем ошибку, но не завершаем программу
		}

		if _, err := db.Exec(query); err != nil {
			log.Printf("Ошибка выполнения SQL-запроса из файла %s: %v\nЗапрос: %s", file, err, query)
			return err
		}
	}

	return nil
}
