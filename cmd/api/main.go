package main

import (
	"bsnack/cmd/middleware"
	"bsnack/config"
	"bsnack/internal/handler/http"
	"bsnack/internal/repository/postgres"
	"bsnack/internal/repository/redis"
	"bsnack/internal/service"
	"bsnack/pkg/database"
	"bsnack/pkg/logger"
	"log"
	netHttp "net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.AppEnv == "production")
	logger.Info("Starting BSNACK API",
		"port", cfg.AppPort,
		"env", cfg.AppEnv,
	)

	db, err := database.NewPostgresDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient(cfg.RedisHost, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	prodRepo := postgres.NewProductRepo(db)
	custRepo := postgres.NewCustomerRepo(db)
	transRepo := postgres.NewTransactionRepo(db)
	cacheRepo := redis.NewRedisRepo(rdb)

	prodSvc := service.NewProductService(prodRepo)
	transSvc := service.NewTransactionService(prodRepo, custRepo, transRepo, cacheRepo)
	custSvc := service.NewCustomerService(custRepo)

	handler := http.NewHandler(prodSvc, transSvc, custSvc)

	mux := netHttp.NewServeMux()

	mux.HandleFunc("GET /customers", handler.ListCustomers)

	mux.HandleFunc("POST /products", handler.AddProduct)
	mux.HandleFunc("GET /products", handler.GetProducts)

	mux.HandleFunc("POST /transactions", handler.CreateTransaction)
	mux.HandleFunc("GET /transactions", handler.GetReport)
	mux.HandleFunc("POST /redemptions", handler.Redeem)

	loggingMiddleware := middleware.RequestLogger(mux)

	serverAddr := ":" + cfg.AppPort

	server := &netHttp.Server{
		Addr:         serverAddr,
		Handler:      loggingMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Server is running", "port", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
