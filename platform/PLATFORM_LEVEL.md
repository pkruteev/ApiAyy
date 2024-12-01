# ./platform

**Folder with platform-level logic**. This directory contains all the platform-level logic that will build up the actual project, like _setting up the database_ or _cache server instance_ and _storing migrations_.
Этот каталог содержит всю логику на уровне платформы, которая будет использоваться для создания реального проекта, например, для настройки базы данных или экземпляра сервера кэширования и для сохранения миграций.

- `./platform/database` folder with database configuration (by default, PostgreSQL)
- `./platform/migrations` folder with migration files (used with [golang-migrate/migrate](https://github.com/golang-migrate/migrate) tool)
