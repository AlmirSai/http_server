package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"http_server/auth-service/internal/config"
	"http_server/auth-service/internal/domain/models"
	"http_server/auth-service/internal/domain/repository"
	"http_server/auth-service/internal/handler"
	"http_server/auth-service/internal/server"
	"http_server/auth-service/internal/service"
	"http_server/auth-service/pkg/logging"
	"http_server/auth-service/pkg/middleware"
	"http_server/auth-service/pkg/monitoring"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Println("Failed to load config: %v\n", err)
	}

	// Initialize logger
	logConfig := &logging.Config{
		Environment: os.Getenv("APP_ENV"),
		LogLevel:    cfg.Logging.Level,
		FilePath:    cfg.Logging.Output,
		MaxSize:     10, // 10MB
		MaxBackups:  5,
		MaxAge:      30, // 30 days
		Compress:    true,
	}

	logger, err := logging.NewLogger(logConfig)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func(logger *logging.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Errorf("Failed to sync logger: %v\n", err)
		}
	}(logger)

	// Initialize metrics
	metrics := monitoring.NewMetrics(cfg.Metrics.ServiceName)

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Auto-migrate database schemas
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}); err != nil {
		logger.Fatal("Failed to auto-migrate database", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, []byte(cfg.JWT.SecretKey), logger)

	// Initialize handlers and middleware
	authHandler := handler.NewAuthHandler(authService, logger, metrics)
	authMiddleware := middleware.NewAuthMiddleware(authService, logger, metrics)

	// Initialize server
	srv := server.NewServer(cfg, authHandler, authMiddleware)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal("Failed to start server", err)
		}
	}()

	logger.Info("Server started successfully")

	// Wait for interrupt signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown completed")
}
