package db

import (
	"echo-server/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(config config.DatabaseConfig) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}
	log.Printf("Database connected - %s:%d/%s/%s", config.Host, config.Port, config.DBName, config.User)
	return database
}
