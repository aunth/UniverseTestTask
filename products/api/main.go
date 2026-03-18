package main

import (
	"context"
	"log"
	"net/http"

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
	log.Println("Starting database migrations...")
	if err := postgres.RunMigrations(cfg.DatabaseURL, cfg.MigrationPath); err != nil && err.Error() != "no change" {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	dbPool, err := postgres.InitDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL connection established")

	return dbPool
}

func setupBroker(cfg *config.Config) *broker.RabbitMQBroker {
	msgBroker, err := broker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v\n", err)
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
	log.Printf("Products Service started on port %s\n", port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
