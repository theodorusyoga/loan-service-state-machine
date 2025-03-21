package migrations

import (
	"log"

	migrations_models "github.com/theodorusyoga/loan-service-state-machine/migrations/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func Migrate(db *gorm.DB) error {
	log.Println("Running migrations")
	err := db.AutoMigrate(&migrations_models.Loan{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Migrations completed successfully")
	return err
}
