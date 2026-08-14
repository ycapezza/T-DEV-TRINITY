package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"t_dev_700/internal/usecase"
	"t_dev_700/pkg/paypal"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	paymentService *usecase.PaymentService
	paypalClient   *paypal.Client
	webhookID      string
}

func NewWebhookHandler(paymentService *usecase.PaymentService, paypalClient *paypal.Client, webhookID string) *WebhookHandler {
	return &WebhookHandler{
		paymentService: paymentService,
		paypalClient:   paypalClient,
		webhookID:      webhookID,
	}
}

func (h *WebhookHandler) HandlePayPalWebhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	headers := map[string]string{
		"PAYPAL-AUTH-ALGO":       c.GetHeader("PAYPAL-AUTH-ALGO"),
		"PAYPAL-CERT-URL":        c.GetHeader("PAYPAL-CERT-URL"),
		"PAYPAL-TRANSMISSION-ID": c.GetHeader("PAYPAL-TRANSMISSION-ID"),
		"PAYPAL-TRANSMISSION-SIG": c.GetHeader("PAYPAL-TRANSMISSION-SIG"),
		"PAYPAL-TRANSMISSION-TIME": c.GetHeader("PAYPAL-TRANSMISSION-TIME"),
	}

	valid, err := h.paypalClient.VerifyWebhookSignature(headers, body, h.webhookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify webhook signature"})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook event"})
		return
	}

	switch event.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		if err := h.paymentService.HandlePaymentCompleted(c, event.Resource.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment completion"})
			return
		}
	case "PAYMENT.CAPTURE.DENIED":
		if err := h.paymentService.HandlePaymentDenied(c, event.Resource.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment denial"})
			return
		}
	case "PAYMENT.CAPTURE.REFUNDED":
		if err := h.paymentService.HandlePaymentRefunded(c, event.Resource.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment refund"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}