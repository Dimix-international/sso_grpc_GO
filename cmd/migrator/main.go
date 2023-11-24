package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" //драйвер для sqlite 3
	_ "github.com/golang-migrate/migrate/v4/source/file"      //драйвер для получений миграций из файлов
)

//go run .\cmd\migrator\ --storage-path=./storage/sso.db --migrations-path=./migrations - команда запуска
func main() {
	var storagePath, migrationsPath, migrationsTable string
	//storagePath - путь к бд - где будет храниться
	//migrationsPath - где находятся файлы миграции
	//migrationsTable - указываем имя таблицы для какой сохраняем (для тестов функциональных пригодится) - необязательно

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	//создаем экземпляр мигратора
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable), //?x-migrations-table=%s - необязат параметр для migrationsTable
	)
	if err != nil {
		panic(err)
	}

	//выполняем миграцию
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			//все миграции применимы
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}