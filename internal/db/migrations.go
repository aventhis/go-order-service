package db

import (
	"log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)


func RunMigrations(db *sqlx.DB, migrationsPath string)  {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %v", err)
	}
	
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", 
		driver,
	)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration failed: %v", err)
    }
    log.Println("Migrations applied successfully")
}
