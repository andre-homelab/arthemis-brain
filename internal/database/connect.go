package database

import (
	"fmt"
	"time"

	"arthemis-brain/internal/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	host := env.GetEnv("DB_HOST", "arthemis-brain-postgres")
	port := env.GetEnv("DB_PORT", "5432")
	user := env.GetEnv("DB_USER", "app_user")
	password := env.GetEnv("DB_PASSWORD", "app_password")
	databaseName := env.GetEnv("DB_NAME", "app_db")
	// dev := env.GetEnv("DEVELOPMENT", "true")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		databaseName,
	)

	db, err := gorm.Open(
		postgres.Open(connectionString),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return db, nil
}
