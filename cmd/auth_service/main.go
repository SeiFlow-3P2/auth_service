package main

import (
	"github.com/SeiFlow-3P2/auth_service/internal/app"
	"github.com/SeiFlow-3P2/auth_service/internal/service"
	"log/slog"
)

const (
	configPath = ".env" //читается в /internal/app/app.go  !!!!
)

func main() {
	authApp := app.NewApp(configPath)
	if authApp == nil {
		panic("app is nil")
	}

	auth := service.Auth{authApp}
	grpcServer := app.NewGRPCApp(slog.Default(), &auth, configPath)

	err := grpcServer.Run()
	if err != nil {
		panic(err)
	}
}
