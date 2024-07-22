package db

import (
	"fmt"
	"onboarding_test/internal/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s",
		config.Envs.DB,
		config.Envs.DBName,
		config.Envs.DBPassword,
		config.Envs.Host,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
		return nil, err
	}
	fmt.Println("Connected to database successfully")
	return db, nil
}