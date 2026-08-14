package http

import (
	"t_dev_700/internal/delivery/http/handler"
	"t_dev_700/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	invoiceHandler *handler.InvoiceHandler,
	reportHandler *handler.ReportHandler,
	authMiddleware *middleware.AuthMiddleware,
	paymentHandler *handler.PaymentHandler,
  webhookHandler *handler.WebhookHandler,
) *gin.Engine {
	r := gin.Default()

	middleware.SetupSecurityMiddleware(r)

	csrfMiddleware := middleware.NewCSRFMiddleware()

	r.POST("/webhooks/paypal", webhookHandler.HandlePayPalWebhook)
	r.GET("/payments/capture", paymentHandler.CapturePayment)
  r.GET("/payments/cancel", paymentHandler.CancelPayment)

	auth := r.Group("/auth")
	{
		auth.POST("/register", csrfMiddleware.ValidateToken(), authHandler.Register)
		auth.POST("/login/admin", csrfMiddleware.ValidateToken(), authHandler.LoginAdmin)
		auth.POST("/login/mobile", csrfMiddleware.ValidateToken(), authHandler.LoginMobile)
		auth.GET("/csrf-token", csrfMiddleware.GenerateToken())
	}

	api := r.Group("/api", authMiddleware.AuthRequired())
	{
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)
		api.DELETE("/profile", userHandler.DeleteProfile)
		api.GET("/invoices/me", invoiceHandler.GetUserInvoices)
		api.POST("/payments", paymentHandler.CreatePayment)
		
		admin := api.Group("/admin", middleware.AdminRequired())
		{
			admin.GET("/users", userHandler.List)
			admin.POST("/users", userHandler.Create)
			admin.PUT("/users/:id", userHandler.Update)
			admin.DELETE("/users/:id", userHandler.Delete)

			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)
			admin.POST("/products/barcode", productHandler.CreateFromBarcode)
			admin.GET("/products", productHandler.List)
			admin.GET("/products/search", productHandler.Search)
			admin.GET("/products/:id", productHandler.GetByID)

			admin.GET("/invoices", invoiceHandler.List)
			admin.PUT("/invoices/:id", invoiceHandler.Update)
			admin.DELETE("/invoices/:id", invoiceHandler.Delete)

			admin.GET("/reports/sales", reportHandler.GetSalesReport)
			admin.GET("/reports/top-products", reportHandler.GetTopProducts)
			admin.GET("/reports/categories", reportHandler.GetCategoryPerformance)
			admin.GET("/reports/stock-alerts", reportHandler.GetStockAlerts)
			admin.GET("/reports/sales-evolution", reportHandler.GetSalesEvolution)
		}
	}

	return r
}
