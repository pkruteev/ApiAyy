package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
)

// loadSQL загружает SQL-запрос из файла
func loadSQL(filename string) (string, error) {
	// Получаем путь к текущему файлу (create_tables.go)
	_, currentFilePath, _, _ := runtime.Caller(0) // Получаем путь к текущему файлу
	currentDir := filepath.Dir(currentFilePath)   // Получаем директорию текущего файла

	// Получаем путь к папке platform (на один уровень выше database)
	platformDir := filepath.Join(currentDir, "..")

	// Формируем полный путь к папке migrations
	migrationsDir := filepath.Join(platformDir, "migrations")

	// Формируем полный путь к файлу
	filePath := filepath.Join(migrationsDir, filename)

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
func CreateTables(db *sqlx.DB) error {
	// Установка временной зоны
	if _, err := db.Exec("SET TIMEZONE='Europe/Moscow';"); err != nil {
		return err
	}

	// Список SQL-файлов для создания таблиц
	sqlFiles := []string{
		"000003_create_companies.up.sql",
		"000004_create_r_schets.up.sql",
		"000005_create_objects.up.sql",
		"000006_create_contracts.up.sql",
		"000007_create_statements.up.sql",
		"000008_create_expenses.up.sql",
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
