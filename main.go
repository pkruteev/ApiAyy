package main

import (
	"ApiAyy/pkg/configs"
	"ApiAyy/pkg/middleware"
	"ApiAyy/pkg/routes"
	"ApiAyy/pkg/utils"
	"ApiAyy/platform/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Определяем конфигурацию Fiber.
	config := configs.FiberConfig()

	// Создаем новое приложение Fiber с конфигурацией.
	app := fiber.New(config)

	// Подключаем middleware.
	middleware.FiberMiddleware(app)

	// Подключение к стартовой БД postgres.
	db_postgres, err := database.DBConnection("postgres")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Создание основной БД
	err = database.CreateDB("main", db_postgres)
	if err != nil {
		log.Fatalf("Ошибка при создании базы данных 'main': %v", err)
	}

	// fmt.Println("Пул соединений успешно создан.")

	// Подключение к базе данных main.
	db_main, err := database.DBConnection("main")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Вызов функции для создания таблицы Users, если не созданы.
	database.CreateStartTables(db_main)

	// Регистрация маршрутов.
	routes.PublicRoutes(app)  // Регистрация публичных маршрутов.
	routes.PrivateRoutes(app) // Регистрация приватных маршрутов.
	routes.NotFoundRoute(app) // Регистрация маршрута для ошибки 404.

	// Запускаем сервер (с корректным завершением).
	utils.StartServer(app)
}
