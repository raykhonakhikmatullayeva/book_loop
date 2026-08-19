package main

import (
	"book_loop/internal/router"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "book_loop/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	shutdownTimeout = 5 * time.Second
)

// @title						AuctionHouse API
// @version					1.0
// @host						localhost:9999
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	logLevel := getLoggerLevel(os.Getenv("LOG_LEVEL"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found", "err", err)
	}
	ctx := context.Background()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		slog.Error("DB_DSN is not set")
		return errors.New("DB_DSN is not set")
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("connect database", "err", err)
		return fmt.Errorf("connect database: %w", err)
	}
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	engine, err := router.New(ctx, pool)
	if err != nil {
		logger.Error("create router", "err", err)
		return fmt.Errorf("create router: %w", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	addr := ":" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: shutdownTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	logger.Info("server started", "port", port)

	<-sigCtx.Done()
	logger.Info("shutdown signal received")

	shCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shCtx); err != nil {
		logger.Error("shutdown", "err", err)
	} else {
		logger.Info("server stopped")
	}

	pool.Close()
	logger.Info("shutdown")
	return nil
}

func getLoggerLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
