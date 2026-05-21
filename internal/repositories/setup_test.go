package repositories

import (
	"os"
	"testing"

	config "realworld-gin/internal/config/database"
	"realworld-gin/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestMain bootstraps a shared in-memory SQLite database for all repository tests.
func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to open sqlite in-memory DB: " + err.Error())
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Tag{},
		&models.Comment{},
	); err != nil {
		panic("failed to auto-migrate: " + err.Error())
	}

	config.DB = db
	os.Exit(m.Run())
}

// cleanTables truncates all application tables before each test to ensure isolation.
// SQLite does not enforce FK constraints by default, so order does not matter.
func cleanTables(t *testing.T) {
	t.Helper()
	tables := []string{
		"user_favorites",
		"user_follows",
		"article_tags",
		"comments",
		"articles",
		"tags",
		"users",
	}
	for _, table := range tables {
		if err := config.DB.Exec("DELETE FROM " + table).Error; err != nil {
			t.Logf("cleanup warning for table %s: %v", table, err)
		}
	}
}
