package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Error: Not config DB_DSN in .env file")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error: Failed to connect to database:", err)
	}

	log.Println("✅ Connect to MySQL success!")
	DB = database
}
