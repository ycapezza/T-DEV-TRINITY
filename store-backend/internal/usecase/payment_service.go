// internal/usecase/payment_service.go
package usecase

import (
	"context"
	"errors"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/pkg/paypal"
)

type PaymentService struct {
	invoiceRepo  repository.InvoiceRepository
	productRepo  repository.ProductRepository
	paypalClient *paypal.Client
	baseURL      string
}

func NewPaymentService(
	invoiceRepo repository.InvoiceRepository,
	productRepo repository.ProductRepository,
	paypalClient *paypal.Client,
	baseURL string,
) *PaymentService {
	return &PaymentService{
		invoiceRepo:  invoiceRepo,
		productRepo:  productRepo,
		paypalClient: paypalClient,
		baseURL:      baseURL,
	}
}

type CreatePaymentResponse struct {
	OrderID     string `json:"order_id"`
	ApproveURL  string `json:"approve_url"`
	InvoiceID   uint   `json:"invoice_id"`
}

func (s *PaymentService) CreatePayment(ctx context.Context, userID uint, items []struct {
	ProductID uint
	Quantity  int
}) (*CreatePaymentResponse, error) {
	var total float64
	var orderItems []models.OrderItem

	for _, item := range items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.StockQuantity < item.Quantity {
			return nil, errors.New("insufficient stock")
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		total += product.Price * float64(item.Quantity)

		if err := s.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
			return nil, err
		}
	}

	invoice := &models.Invoice{
		UserID: userID,
		Total:  total,
		Status: "pending",
	}

	if err := s.invoiceRepo.Create(ctx, invoice, orderItems); err != nil {
		return nil, err
	}

	successURL := s.baseURL + "/payments/capture"
	cancelURL := s.baseURL + "/payments/cancel"

	orderID, approveURL, err := s.paypalClient.CreateOrder(
		total,
		invoice.ID,
		successURL,
		cancelURL,
	)
	if err != nil {
		return nil, err
	}

	invoice.PaypalID = orderID
	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		OrderID:     orderID,
		ApproveURL:  approveURL,
		InvoiceID:   invoice.ID,
	}, nil
}

func (s *PaymentService) CapturePayment(ctx context.Context, orderID string) error {
	invoice, err := s.invoiceRepo.GetByPayPalID(ctx, orderID)
	if err != nil {
		return err
	}

	captureID, err := s.paypalClient.CaptureOrder(orderID)
	if err != nil {
		return err
	}

	invoice.Status = "paid"
	invoice.PaypalID = captureID
	return s.invoiceRepo.Update(ctx, invoice)
}

func (s *PaymentService) CancelPayment(ctx context.Context, orderID string) error {
	invoice, err := s.invoiceRepo.GetByPayPalID(ctx, orderID)
	if err != nil {
		return err
	}

	_, items, err := s.invoiceRepo.GetWithItems(ctx, invoice.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}

	invoice.Status = "cancelled"
	return s.invoiceRepo.Update(ctx, invoice)
}

func (s *PaymentService) HandlePaymentCompleted(ctx context.Context, captureID string) error {
	invoice, err := s.invoiceRepo.GetByPayPalID(ctx, captureID)
	if err != nil {
		return err
	}

	invoice.Status = "paid"
	return s.invoiceRepo.Update(ctx, invoice)
}

func (s *PaymentService) HandlePaymentDenied(ctx context.Context, captureID string) error {
	invoice, err := s.invoiceRepo.GetByPayPalID(ctx, captureID)
	if err != nil {
		return err
	}

	_, items, err := s.invoiceRepo.GetWithItems(ctx, invoice.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}

	invoice.Status = "denied"
	return s.invoiceRepo.Update(ctx, invoice)
}

func (s *PaymentService) HandlePaymentRefunded(ctx context.Context, captureID string) error {
	invoice, err := s.invoiceRepo.GetByPayPalID(ctx, captureID)
	if err != nil {
		return err
	}

	invoice.Status = "refunded"
	return s.invoiceRepo.Update(ctx, invoice)
}