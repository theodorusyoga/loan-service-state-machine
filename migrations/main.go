package migrations

import (
	"log"

	"github.com/theodorusyoga/loan-service-state-machine/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunMigrations() {
	cfg, err := config.Load("config/config.yaml")

	if err != nil {
		log.Fatalf("failed to load config for migrations: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
