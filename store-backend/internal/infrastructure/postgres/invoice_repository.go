package postgres

import (
	"context"
	"fmt"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type invoiceRepository struct {
   db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) repository.InvoiceRepository {
   return &invoiceRepository{
       db: db,
   }
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *models.Invoice, items []models.OrderItem) error {
   return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
       if err := tx.Create(invoice).Error; err != nil {
           return err
       }

       for i := range items {
           items[i].InvoiceID = invoice.ID
       }

       return tx.Create(&items).Error
   })
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uint) (*models.Invoice, error) {
   var invoice models.Invoice
   err := r.db.WithContext(ctx).
       Preload("User").
       Preload("Items").
       Preload("Items.Product").
       First(&invoice, id).Error
   return &invoice, err
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
   return r.db.WithContext(ctx).Save(invoice).Error
}

func (r *invoiceRepository) Delete(ctx context.Context, id uint) error {
   return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
       if err := tx.Where("invoice_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
           return err
       }
       return tx.Delete(&models.Invoice{}, id).Error
   })
}

func (r *invoiceRepository) List(ctx context.Context) ([]models.Invoice, error) {
   var invoices []models.Invoice
   err := r.db.WithContext(ctx).
       Preload("User").
       Preload("Items").
       Preload("Items.Product").
       Find(&invoices).Error
   return invoices, err
}

func (r *invoiceRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Invoice, error) {
   var invoices []models.Invoice
   err := r.db.WithContext(ctx).
       Where("user_id = ?", userID).
       Preload("Items").
       Preload("Items.Product").
       Find(&invoices).Error
   return invoices, err
}

func (r *invoiceRepository) GetSalesBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]models.Invoice, error) {
    var invoices []models.Invoice
    err := r.db.WithContext(ctx).
        Where("created_at BETWEEN ? AND ?", startDate, endDate).
        Find(&invoices).Error
    return invoices, err
}

func (r *invoiceRepository) GetTopSellingProducts(ctx context.Context, limit int) ([]models.ProductSalesReport, error) {
    var reports []models.ProductSalesReport
    
    err := r.db.WithContext(ctx).
        Table("order_items").
        Select(`
            product_id,
            products.name as product_name,
            SUM(quantity) as quantity_sold,
            SUM(order_items.price * quantity) as total_revenue
        `).
        Joins("JOIN products ON products.id = order_items.product_id").
        Group("product_id, products.name").
        Order("quantity_sold DESC").
        Limit(limit).
        Scan(&reports).Error

    return reports, err
}

func (r *invoiceRepository) GetSalesByCategory(ctx context.Context) ([]models.CategoryReport, error) {
    var reports []models.CategoryReport
    
    err := r.db.WithContext(ctx).
        Table("order_items").
        Select(`
            COALESCE(products.category, 'Uncategorized') as category,
            SUM(order_items.price * order_items.quantity) as total_sales,
            COUNT(DISTINCT order_items.invoice_id) as order_count
        `).
        Joins("JOIN products ON products.id = order_items.product_id").
        Group("products.category").
        Scan(&reports).Error

    if err != nil {
        return nil, err
    }
    
    if reports == nil {
        reports = []models.CategoryReport{} 
    }

    return reports, nil
}

func (r *invoiceRepository) GetSalesEvolution(ctx context.Context, period string) ([]models.SalesEvolution, error) {
    var reports []models.SalesEvolution
    
    groupBy := "DATE(created_at)"
    if period == "monthly" {
        groupBy = "DATE_TRUNC('month', created_at)"
    }

    err := r.db.WithContext(ctx).
        Table("invoices").
        Select(fmt.Sprintf(`
            %s as period,
            SUM(total) as revenue,
            COUNT(*) as order_count
        `, groupBy)).
        Group(groupBy).
        Order(groupBy).
        Scan(&reports).Error

    return reports, err
}

func (r *invoiceRepository) GetByPayPalID(ctx context.Context, paypalID string) (*models.Invoice, error) {
    var invoice models.Invoice
    err := r.db.WithContext(ctx).Where("paypal_id = ?", paypalID).First(&invoice).Error
    return &invoice, err
}

func (r *invoiceRepository) GetWithItems(ctx context.Context, id uint) (*models.Invoice, []models.OrderItem, error) {
    var invoice models.Invoice
    err := r.db.WithContext(ctx).First(&invoice, id).Error
    if err != nil {
        return nil, nil, err
    }

    var items []models.OrderItem
    err = r.db.WithContext(ctx).Where("invoice_id = ?", id).Find(&items).Error
    if err != nil {
        return nil, nil, err
    }

    return &invoice, items, nil
}