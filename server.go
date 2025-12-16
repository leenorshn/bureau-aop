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

	// Connect to MongoDB with optimized connection pool settings for Cloud Run
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(50).                         // Maximum number of connections in the pool
		SetMinPoolSize(5).                          // Minimum number of connections to maintain
		SetMaxConnIdleTime(30 * time.Second).       // Close connections after 30s of inactivity
		SetConnectTimeout(10 * time.Second).        // Timeout for initial connection
		SetServerSelectionTimeout(5 * time.Second). // Timeout for server selection
		SetSocketTimeout(30 * time.Second).         // Timeout for socket operations
		SetHeartbeatInterval(10 * time.Second)      // Heartbeat interval for connection monitoring

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Ping the primary to verify connection with timeout
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}
	logger.Info("Connected to MongoDB!",
		zap.String("maxPoolSize", "50"),
		zap.String("minPoolSize", "5"))

	db := client.Database(cfg.MongoDBName)

	// Initialize repositories
	productRepo := store.NewProductRepository(db)
	clientRepo := store.NewClientRepository(db)
	saleRepo := store.NewSaleRepository(db, logger)
	paymentRepo := store.NewPaymentRepository(db)
	commissionRepo := store.NewCommissionRepository(db)
	adminRepo := store.NewAdminRepository(db)
	caisseRepo := store.NewCaisseRepository(db)
	binaryCappingRepo := store.NewBinaryCappingRepository(db)

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
	caisseService := service.NewCaisseService(caisseRepo, logger)
	
	// Initialize Binary Commission Service with new algorithm
	binaryConfig := models.BinaryConfig{
		CycleValue:         cfg.BinaryCycleValue,
		DailyCycleLimit:    cfg.BinaryDailyCycleLimit,
		WeeklyCycleLimit:   cfg.BinaryWeeklyCycleLimit,
		MinVolumePerLeg:    cfg.BinaryMinVolumePerLeg,
		RequireDirectLeft:  true,
		RequireDirectRight: true,
	}
	binaryCommissionService := service.NewBinaryCommissionService(
		clientRepo,
		commissionRepo,
		saleRepo,
		binaryCappingRepo,
		logger,
		binaryConfig,
	)

	// Initialize GraphQL resolver
	resolver := graph.NewResolver(
		productService,
		clientService,
		saleService,
		paymentService,
		commissionService,
		authService,
		adminService,
		caisseService,
		binaryCommissionService,
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

	// Configure HTTP server with timeouts optimized for Cloud Run
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           nil,              // Handler is set by http.Handle
		ReadTimeout:       15 * time.Second, // Maximum duration for reading the entire request
		WriteTimeout:      15 * time.Second, // Maximum duration before timing out writes
		IdleTimeout:       60 * time.Second, // Maximum amount of time to wait for the next request
		ReadHeaderTimeout: 5 * time.Second,  // Amount of time allowed to read request headers
	}

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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server exited")
}
