package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"catalog-product/internal/broker"
	"catalog-product/internal/config"
	"catalog-product/internal/handler"
	"catalog-product/internal/repository"
	"catalog-product/internal/service"
	"catalog-product/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	dbPool := setupDatabase(ctx, cfg)
	defer dbPool.Close()

	msgBroker := setupBroker(cfg)
	if msgBroker != nil {
		defer msgBroker.Close()
	}

	router := setupRouter(dbPool, msgBroker)

	runServer(router, cfg.Port)
}

func setupDatabase(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	slog.Info("Starting database migrations...")
	if err := postgres.RunMigrations(cfg.DatabaseURL, cfg.MigrationPath); err != nil && err.Error() != "no change" {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Migrations applied successfully")

	dbPool, err := postgres.InitDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("PostgreSQL connection established")

	return dbPool
}

func setupBroker(cfg *config.Config) *broker.RabbitMQBroker {
	msgBroker, err := broker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return nil
	}
	return msgBroker
}

func setupRouter(dbPool *pgxpool.Pool, msgBroker *broker.RabbitMQBroker) *gin.Engine {
	productRepo := repository.NewProductRepository(dbPool)
	productSvc := service.NewProductService(productRepo, msgBroker)
	productHandler := handler.NewProductHandler(productSvc)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "connected"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	productHandler.RegisterRoutes(router)

	return router
}

func runServer(router *gin.Engine, port string) {
	addr := ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		slog.Info("Products Service started", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
