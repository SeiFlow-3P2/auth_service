package app

import (
	"context"
	"fmt"
	"github.com/SeiFlow-3P2/auth_service/internal/domain"
	"github.com/SeiFlow-3P2/auth_service/pkg/authOrm"
	"github.com/SeiFlow-3P2/auth_service/pkg/authRedis"
	"github.com/SeiFlow-3P2/auth_service/pkg/grpc/auth_v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"os"
	"strconv"
	"time"
)

func NewApp(cfgPath string) *domain.App {

	//Загрузка переменных окружения
	if err := godotenv.Load(cfgPath); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// Инициализация БД
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbMode := os.Getenv("DB_SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbUser, dbName, dbPass, dbMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error opening DB: %v", err))
	}

	authDB := &authOrm.AuthOrm{*db}
	err = authDB.MigrateDB()

	if err != nil {
		panic(fmt.Sprintf("Error migrating DB: %v", err))
	}
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

	if redis == nil {
		panic("cant create redis client")
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		panic("cant parse app url")
	}
	configs := make(map[string]*oauth2.Config)

	configs["github"] = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  appUrl + "/callback/github",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	configs["google"] = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  appUrl + "callback/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return &domain.App{
		AuthDB: authDB,
		Casher: redis,
		Settings: &domain.AppSettings{
			Secret:     secret,
			RefreshTTL: refreshTTL,
			AccessTTL:  accessTTL,
		},
		Logger: slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		OauthConfigs: configs,
	}
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// NewGRPCApp creates new gRPC server app.
func NewGRPCApp(
	log *slog.Logger,
	authService auth_v1.Auth,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	auth_v1.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
