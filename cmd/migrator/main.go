package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Novochenko/sso/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	// flag.Parse()
	cfg := config.MustLoad()
	var DBUrl string
	log.Println(cfg.Env)
	if cfg.Env == "local" {
		DBUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", cfg.StoragePath.User, cfg.StoragePath.Password, cfg.StoragePath.Host, cfg.StoragePath.DBName)
	} else {
		DBUrl = cfg.StoragePath.FullName
	}
	db, err := sql.Open("mysql", DBUrl)
	if err != nil {
		panic(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}
