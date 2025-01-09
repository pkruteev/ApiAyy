package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// PostgreNewSQLConnection создает подключение к базе данных PostgreSQL
func DBConnection(db_name string) (*sqlx.DB, error) {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки файла .env, %w", err)
	}

	// Определяем параметры подключения к базе данных.
	maxConn, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании DB_MAX_CONNECTIONS, %w", err)
	}

	maxIdleConn, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании DB_MAX_IDLE_CONNECTIONS, %w", err)
	}

	maxLifetimeConn, err := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании DB_MAX_LIFETIME_CONNECTIONS, %w", err)
	}

	// Получаем настройку из файла ENV
	dbConfig := os.Getenv("DB_NEW_SERVER_URL")

	// Находим имя базы данных
	dbname := fmt.Sprintf("dbname=%s", db_name)

	// Формируем URL для подключения
	dbURL := fmt.Sprintf("%s %s", dbConfig, dbname)

	// Подключаемся к бд
	db_new, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка, не удалось подключиться к базе данных, %w", err)
	}

	// Устанавливаем параметры подключения к базе данных.
	db_new.SetMaxOpenConns(maxConn)                           // по умолчанию 0 (неограниченно)
	db_new.SetMaxIdleConns(maxIdleConn)                       // по умолчанию 2
	db_new.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, соединения повторно используются навсегда

	// Пытаемся отправить запрос ping к базе данных.
	if err := db_new.Ping(); err != nil {
		defer db_new.Close() // закрываем подключение к базе данных
		return nil, fmt.Errorf("ошибка ping, не удалось отправить ping к базе данных, %w", err)
	}

	return db_new, nil
}
