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

func PostgreSQLConnection() (*sqlx.DB, error) {
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

	// Подключаемся к базе данных PostgreSQL.
	db_main, err := sqlx.Connect("pgx", os.Getenv("DB_SERVER_URL"))
	if err != nil {
		return nil, fmt.Errorf("ошибка, не удалось подключиться к базе данных, %w", err)
	}

	// Устанавливаем параметры подключения к базе данных.
	db_main.SetMaxOpenConns(maxConn)                           // по умолчанию 0 (неограниченно)
	db_main.SetMaxIdleConns(maxIdleConn)                       // по умолчанию 2
	db_main.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, соединения повторно используются навсегда

	// Пытаемся отправить запрос ping к базе данных.
	if err := db_main.Ping(); err != nil {
		defer db_main.Close() // закрываем подключение к базе данных
		return nil, fmt.Errorf("ошибка ping, не удалось отправить ping к базе данных, %w", err)
	}

	return db_main, nil
}
