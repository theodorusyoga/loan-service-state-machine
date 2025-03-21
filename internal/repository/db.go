package repository

import (
	"log"

	"github.com/theodorusyoga/loan-service-state-machine/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.Database.URL

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to DB successfully")

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Adjust these values based on your application needs
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// TODO: Add auto migration for GORM
func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}
