package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func CreateBooksTable(db *sqlx.DB) error {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущей директории: %v", err)
	}
	fmt.Println("Текущая директория:", dir)

	// Чтение SQL-запроса из файла

	sqlQuery, err := os.ReadFile("create_tables_book.up.sql")
	if err != nil {
		log.Fatalf("Не удалось прочитать SQL файл: %v", err)
		return err
	}

	// Выполнение SQL-запроса
	_, err = db.Exec(string(sqlQuery))
	return err
}
