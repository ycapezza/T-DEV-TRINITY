package usecase

import (
	"context"
	"errors"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
)

type InvoiceService struct {
    invoiceRepo repository.InvoiceRepository
    productRepo repository.ProductRepository
}

func NewInvoiceService(invoiceRepo repository.InvoiceRepository, productRepo repository.ProductRepository) *InvoiceService {
    return &InvoiceService{
        invoiceRepo: invoiceRepo,
        productRepo: productRepo,
    }
}

type OrderItemInput struct {
    ProductID uint    `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

func (s *InvoiceService) Create(ctx context.Context, input struct {
    UserID uint
    Items  []OrderItemInput
}) (*models.Invoice, error) {
    var total float64
    var orderItems []models.OrderItem

    for _, item := range input.Items {
        product, err := s.productRepo.GetByID(ctx, item.ProductID)
        if err != nil {
            return nil, err
        }
        if product.StockQuantity < item.Quantity {
            return nil, errors.New("insufficient stock")
        }

        total += item.Price * float64(item.Quantity)
        orderItems = append(orderItems, models.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     item.Price,
        })

        if err := s.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
            return nil, err
        }
    }

    invoice := &models.Invoice{
        UserID: input.UserID,
        Total:  total,
        Status: "pending",
    }

    if err := s.invoiceRepo.Create(ctx, invoice, orderItems); err != nil {
        return nil, err
    }

    return invoice, nil
}

func (s *InvoiceService) GetByID(ctx context.Context, id uint) (*models.Invoice, error) {
    return s.invoiceRepo.GetByID(ctx, id)
}

func (s *InvoiceService) Update(ctx context.Context, id uint, input struct {
    Status string
}) (*models.Invoice, error) {
    invoice, err := s.invoiceRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    invoice.Status = input.Status

    if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
        return nil, err
    }

    return invoice, nil
}

func (s *InvoiceService) Delete(ctx context.Context, id uint) error {
    invoice, err := s.invoiceRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    for _, item := range invoice.Items {
        if err := s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
            return err
        }
    }

    return s.invoiceRepo.Delete(ctx, id)
}

func (s *InvoiceService) List(ctx context.Context) ([]models.Invoice, error) {
    return s.invoiceRepo.List(ctx)
}

func (s *InvoiceService) GetByUserID(ctx context.Context, userID uint) ([]models.Invoice, error) {
    return s.invoiceRepo.GetByUserID(ctx, userID)
}

func (s *InvoiceService) GetUserInvoices(ctx context.Context, userID uint) ([]models.Invoice, error) {
    return s.invoiceRepo.GetByUserID(ctx, userID)
}