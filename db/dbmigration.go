package db

import (
	"database/sql"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(conn *sql.DB) {
	log.Println("Database migration start")
	driver, _ := postgres.WithInstance(conn, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		log.Println("Database migration failed : ", err)
	}

	if err := m.Down(); err != nil {
		log.Println("Database migration failed down : ", err)
	}
	if err := m.Up(); err != nil {
		log.Println("Database migration failed up : ", err)
	}
	log.Println("Database migration end")
}
