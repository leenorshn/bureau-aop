package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bureau/services/tree-service/internal/cache"
	"bureau/services/tree-service/internal/config"
	"bureau/services/tree-service/internal/handler"
	"bureau/services/tree-service/internal/service"
	"bureau/services/tree-service/internal/store"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const defaultPort = "8082"

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("Failed to sync logger:", err)
		}
	}()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Ping MongoDB
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}
	logger.Info("Connected to MongoDB")

	db := mongoClient.Database(cfg.MongoDBName)

	// Initialize Redis cache (optional)
	var treeCache cache.TreeCache
	if cfg.RedisURL != "" {
		treeCache = cache.NewRedisCache(cfg.RedisURL, logger)
	} else {
		treeCache = cache.NewMemoryCache(logger)
		logger.Info("Using in-memory cache (Redis not configured)")
	}

	// Initialize repositories
	clientRepo := store.NewClientRepository(db)
	saleRepo := store.NewSaleRepository(db, logger)

	// Initialize services
	treeService := service.NewTreeService(clientRepo, saleRepo, treeCache, logger)

	// Initialize HTTP handler
	treeHandler := handler.NewTreeHandler(treeService, logger)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/v1/tree/", treeHandler.HandleTreeRequest)

	// Start server
	port := os.Getenv("TREE_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("Starting Tree Service", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down Tree Service...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Tree Service exited")
}



