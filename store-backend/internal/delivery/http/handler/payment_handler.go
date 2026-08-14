package handler

import (
	"net/http"
	"t_dev_700/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *usecase.PaymentService
}

func NewPaymentHandler(service *usecase.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

type createPaymentRequest struct {
	Items []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,dive"`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
    var req createPaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, _ := c.Get("user_id")

    items := make([]struct {
        ProductID uint
        Quantity  int
    }, len(req.Items))
    
    for i, item := range req.Items {
        items[i].ProductID = item.ProductID
        items[i].Quantity = item.Quantity
    }

    payment, err := h.service.CreatePayment(c.Request.Context(), userID.(uint), items)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) CapturePayment(c *gin.Context) {
    token := c.Query("token")
    payerID := c.Query("PayerID")
    
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID/token is required"})
        return
    }

    if err := h.service.CapturePayment(c.Request.Context(), token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Payment captured successfully",
        "token": token,
        "payer_id": payerID,
    })
}

func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	if err := h.service.CancelPayment(c.Request.Context(), orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment cancelled successfully"})
}