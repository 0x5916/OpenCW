package databases

import (
	"fmt"
	"log/slog"
	"os"

	"opencw/internal/configs"
	"opencw/internal/models"

	"github.com/google/uuid"
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
		&models.ForumCategory{},
		&models.ForumPost{},
		&models.ForumThread{},
	); err != nil {
		slog.Error("Failed to migrate database", "err", err)
	}

	if err := ensureDefaultForumCategory(db); err != nil {
		slog.Error("Failed to seed default forum category", "err", err)
	}

	if err := ensureForumIndexes(db); err != nil {
		slog.Error("Failed to migrate forum indexes", "err", err)
	}

	if err := ensureUserEmailIndexes(db); err != nil {
		slog.Error("Failed to migrate user email indexes", "err", err)
	}

	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get underlying sql.DB", "err", err)
		os.Exit(1)
	}
	sqlDB.SetMaxOpenConns(configs.App.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(configs.App.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(configs.App.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(configs.App.DBConnMaxIdleTime)

	slog.Info("Database connected and migrated successfully")
}

func ensureForumIndexes(db *gorm.DB) error {
	statements := []string{
		`CREATE INDEX IF NOT EXISTS idx_forum_threads_category_cursor ON forum_thread (category_id, is_pinned DESC, updated_at DESC, id DESC) WHERE deleted_at IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_forum_posts_thread_cursor ON forum_post (thread_id, created_at ASC, id ASC) WHERE deleted_at IS NULL;`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}

func ensureDefaultForumCategory(db *gorm.DB) error {
	const defaultCategoryName = "General"

	var category models.ForumCategory
	result := db.Where("name = ?", defaultCategoryName).Limit(1).Find(&category)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	return db.Create(&models.ForumCategory{
		ID:          uuid.MustParse("9f3f1b2c-4e5d-4f6a-8b7c-1a2b3c4d5e6f"),
		Name:        defaultCategoryName,
		Description: "General discussion about CW and OpenCW.",
	}).Error
}

func ensureUserEmailIndexes(db *gorm.DB) error {
	statements := []string{
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;`,
		`DROP INDEX IF EXISTS idx_users_email;`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_verified_email ON users (email) WHERE email_verified = true AND deleted_at IS NULL;`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}

	return nil
}
