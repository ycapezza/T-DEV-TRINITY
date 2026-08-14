package usecase

import (
	"context"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/pkg/pagination"
	"time"
)

type ReportService struct {
    invoiceRepo repository.InvoiceRepository
    productRepo repository.ProductRepository
}

func NewReportService(
    invoiceRepo repository.InvoiceRepository,
    productRepo repository.ProductRepository,
) *ReportService {
    return &ReportService{
        invoiceRepo: invoiceRepo,
        productRepo: productRepo,
    }
}

func (s *ReportService) GenerateSalesReport(ctx context.Context, startDate, endDate time.Time) (*models.SalesReport, error) {
    invoices, err := s.invoiceRepo.GetSalesBetweenDates(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }

    var totalRevenue float64
    for _, invoice := range invoices {
        totalRevenue += invoice.Total
    }

    averageOrderSize := 0.0
    if len(invoices) > 0 {
        averageOrderSize = totalRevenue / float64(len(invoices))
    }

    return &models.SalesReport{
        TotalRevenue:     totalRevenue,
        TotalOrders:      len(invoices),
        AverageOrderSize: averageOrderSize,
    }, nil
}

func (s *ReportService) GetTopProducts(ctx context.Context, limit int) ([]models.ProductSalesReport, error) {
    return s.invoiceRepo.GetTopSellingProducts(ctx, limit)
}

func (s *ReportService) GetCategoryPerformance(ctx context.Context) ([]models.CategoryReport, error) {
    return s.invoiceRepo.GetSalesByCategory(ctx)
}

func (s *ReportService) GetStockAlerts(ctx context.Context, minimumStock int) ([]models.StockAlert, error) {
    p := &pagination.Pagination{
        Page:     1,
        PageSize: 1000,
    }
    
    products, err := s.productRepo.List(ctx, p, nil)
    if err != nil {
        return nil, err
    }

    var alerts []models.StockAlert
    for _, product := range products {
        if product.StockQuantity < minimumStock {
            alerts = append(alerts, models.StockAlert{
                ProductID:    product.ID,
                ProductName:  product.Name,
                CurrentStock: product.StockQuantity,
                MinimumStock: minimumStock,
            })
        }
    }

    return alerts, nil
}

func (s *ReportService) GetSalesEvolution(ctx context.Context, period string) ([]models.SalesEvolution, error) {
    return s.invoiceRepo.GetSalesEvolution(ctx, period)
}