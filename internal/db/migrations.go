package db

import (
	"fmt"
	"log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)


func RunMigrations(db *sqlx.DB, migrationsPath string)  {
	log.Printf("Running migrations from path: %s", migrationsPath)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %v", err)
	}
	
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", 
		driver,
	)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration failed: %v", err)
    }
    
	// Проверяем версию миграции
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Error getting migration version: %v", err)
	} else {
		log.Printf("Current migration version: %d, Dirty: %v", version, dirty)
	}
    
	log.Println("Migrations applied successfully")
}
