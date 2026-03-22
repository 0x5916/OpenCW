package databases

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"opencw/configs"
	"opencw/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		configs.App.DBHost,
		configs.App.DBUser,
		configs.App.DBPassword,
		configs.App.DBName,
		configs.App.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true, PrepareStmt: true})
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.EmailOTP{},
		&models.RefreshToken{},
		&models.CWSettings{},
		&models.PageSettings{},
		&models.Progress{},
	); err != nil {
		slog.Error("Failed to migrate database", "err", err)
	}

	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get underlying sql.DB", "err", err)
		os.Exit(1)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	slog.Info("Database connected and migrated successfully")
}
