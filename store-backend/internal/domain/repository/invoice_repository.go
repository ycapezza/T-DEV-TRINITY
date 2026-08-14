package repository

import (
	"context"
	"t_dev_700/internal/domain/models"
	"time"
)

type InvoiceRepository interface {
    Create(ctx context.Context, invoice *models.Invoice, items []models.OrderItem) error
    GetByID(ctx context.Context, id uint) (*models.Invoice, error)
    Update(ctx context.Context, invoice *models.Invoice) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context) ([]models.Invoice, error)
    GetByUserID(ctx context.Context, userID uint) ([]models.Invoice, error)
    GetSalesBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]models.Invoice, error)
    GetTopSellingProducts(ctx context.Context, limit int) ([]models.ProductSalesReport, error)
    GetSalesByCategory(ctx context.Context) ([]models.CategoryReport, error)
    GetSalesEvolution(ctx context.Context, period string) ([]models.SalesEvolution, error)
    GetByPayPalID(ctx context.Context, paypalID string) (*models.Invoice, error)
    GetWithItems(ctx context.Context, id uint) (*models.Invoice, []models.OrderItem, error)
}