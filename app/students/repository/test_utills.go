package repository

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"echo-server/app/students/model"
)

func SetupTestDB(t *testing.T, dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	err = db.AutoMigrate(
		&model.Student{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}
