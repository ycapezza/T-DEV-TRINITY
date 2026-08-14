package main

import (
	"flag"
	"fmt"
	"log"
	httpDelivery "t_dev_700/internal/delivery/http"
	"t_dev_700/internal/delivery/http/handler"
	"t_dev_700/internal/delivery/http/middleware"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/infrastructure/postgres"
	"t_dev_700/internal/usecase"
	"t_dev_700/pkg/auth"
	"t_dev_700/pkg/bootstrap"
	"t_dev_700/pkg/config"
	"t_dev_700/pkg/openfoodfacts"
	"t_dev_700/pkg/paypal"
	"t_dev_700/pkg/seed"

	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	shouldSeed := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	cfg := config.LoadConfig()

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(pgDriver.Open(dbURL))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Invoice{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if *shouldSeed {
		log.Println("Seeding database...")
		if err := seed.SeedData(db); err != nil {
			log.Fatal("Failed to seed database:", err)
		}
		log.Println("Database seeded successfully!")
	}

	jwtAuth := auth.NewJWTAuth(cfg.JWTSecretKey)

	authMiddleware := middleware.NewAuthMiddleware(jwtAuth)

	offClient := openfoodfacts.NewClient()

	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	invoiceRepo := postgres.NewInvoiceRepository(db)

	authService := usecase.NewAuthService(userRepo, jwtAuth)
	userService := usecase.NewUserService(userRepo)
	productService := usecase.NewProductService(productRepo, offClient)
	invoiceService := usecase.NewInvoiceService(invoiceRepo, productRepo)
	reportService := usecase.NewReportService(invoiceRepo, productRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	reportHandler := handler.NewReportHandler(reportService)
	

	paypalClient := paypal.NewClient(cfg.PaypalClientID, cfg.PaypalClientSecret, true)
	// ICI v
	paymentService := usecase.NewPaymentService(invoiceRepo, productRepo, paypalClient, "https://4466-46-193-0-57.ngrok-free.app")
	paymentHandler := handler.NewPaymentHandler(paymentService)
	webhookHandler := handler.NewWebhookHandler(paymentService, paypalClient, cfg.PaypalWebhookID)

	if err := bootstrap.CreateInitialAdmin(userRepo); err != nil {
		log.Printf("Failed to create initial admin: %v", err)
	}

	router := httpDelivery.SetupRouter(
		authHandler,
		userHandler,
		productHandler,
		invoiceHandler,
		reportHandler,
		authMiddleware,
		paymentHandler,
    webhookHandler,
	)

	log.Printf("Server starting on :%s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
