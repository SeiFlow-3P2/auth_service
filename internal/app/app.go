package app

import (
	"fmt"
	"github.com/SeiFlow-3P2/auth_service/internal/domain"
	"github.com/SeiFlow-3P2/auth_service/pkg/authOrm"
	"github.com/SeiFlow-3P2/auth_service/pkg/authRedis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func App(cfgPath string) *domain.App {

	//Загрузка переменных окружения
	if err := godotenv.Load(cfgPath); err != nil {
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

	refreshTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		panic("cant parse refresh token ttl")
	}
	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		panic("cant parse access token ttl")
	}

	rdHost := os.Getenv("RD_HOST")
	rdPass := os.Getenv("RD_PASSWORD")
	rdID, err := strconv.Atoi(os.Getenv("RD_ID"))
	if err != nil {
		panic("cant parse redis id")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		panic("cant parse secret")
	}

	redis := authRedis.NewRedisClient(rdHost, rdPass, rdID, refreshTTL)

	return &domain.App{
		AuthDB: &authOrm.AuthOrm{*db},
		Casher: redis,
		Settings: &domain.AppSettings{
			Secret:     secret,
			RefreshTTL: refreshTTL,
			AccessTTL:  accessTTL,
		},
		Logger: slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		),
		//TODO init grpc
	}
}

func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		panic(fmt.Sprintf("cant migrate db: %v", err))
	}
}
