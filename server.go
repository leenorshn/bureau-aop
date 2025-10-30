package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bureau/graph"
	"bureau/internal/auth"
	"bureau/internal/config"
	"bureau/internal/models"
	"bureau/internal/service"
	"bureau/internal/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const defaultPort = "8080"

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
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Ping the primary to verify connection
	if err := client.Ping(context.Background(), nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}
	logger.Info("Connected to MongoDB!")

	db := client.Database(cfg.MongoDBName)

	// Initialize repositories
	productRepo := store.NewProductRepository(db)
	clientRepo := store.NewClientRepository(db)
	saleRepo := store.NewSaleRepository(db, logger)
	paymentRepo := store.NewPaymentRepository(db)
	commissionRepo := store.NewCommissionRepository(db)
	adminRepo := store.NewAdminRepository(db)

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg, logger)

	// Initialize services
	productService := service.NewProductService(productRepo, logger)
	clientService := service.NewClientService(clientRepo, saleRepo, commissionRepo, logger, cfg.BinaryThreshold, cfg.BinaryCommissionRate, cfg.DefaultProductPrice)
	saleService := service.NewSaleService(saleRepo, logger)
	paymentService := service.NewPaymentService(paymentRepo, logger)
	commissionService := service.NewCommissionService(commissionRepo, clientRepo, logger, cfg.BinaryCommissionRate, cfg.BinaryThreshold)
	adminService := service.NewAdminService(adminRepo, clientRepo, productRepo, saleRepo, commissionRepo, logger)
	authService := service.NewAuthService(adminRepo, jwtService, logger)

	// Initialize GraphQL resolver
	resolver := graph.NewResolver(
		productService,
		clientService,
		saleService,
		paymentService,
		commissionService,
		authService,
		adminService,
	)

	// Create GraphQL handler
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Create authentication middleware
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token := authHeader[7:]

				// Validate token and get claims
				claims, err := jwtService.ValidateAccessToken(token)
				if err == nil && claims != nil {
					// Create a mock user from claims (in a real app, you'd fetch from DB)
					adminID, _ := primitive.ObjectIDFromHex(claims.AdminID)
					user := &models.Admin{
						ID:    adminID,
						Email: claims.Email,
						Role:  claims.Role,
					}

					// Add user to context
					ctx := context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				}
			}

			next.ServeHTTP(w, r)
		})
	}

	// Setup routes
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authMiddleware(srv))

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	logger.Info("Starting server", zap.String("port", port))
	server := &http.Server{Addr: ":" + port, Handler: nil} // Handler is set by http.Handle

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server exited")
}
