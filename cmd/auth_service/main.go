package main

import (
	"fmt"
	dbORM "github.com/SeiFlow-3P2/auth_service/pkg/gorm"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {

	//Загрузка переменных окружения
	if err := godotenv.Load("../../configs/.env"); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// Инициализация БД
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_Name")
	dbPass := os.Getenv("DB_PASSWORD")
	dbMode := os.Getenv("DB_SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbUser, dbName, dbPass, dbMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Error opening DB: %v", err))
	}
	migrateDB(db)
}

func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&dbORM.User{})
	if err != nil {
		panic(fmt.Sprintf("cant migrate db: %v", err))
	}
}
