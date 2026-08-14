package handler

import (
	"net/http"
	"strconv"
	"t_dev_700/internal/usecase"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
   service *usecase.InvoiceService
}

func NewInvoiceHandler(service *usecase.InvoiceService) *InvoiceHandler {
   return &InvoiceHandler{
       service: service,
   }
}

type createInvoiceRequest struct {
   UserID uint                  `json:"user_id" binding:"required"`
   Items  []createInvoiceItem   `json:"items" binding:"required,dive"`
}

type createInvoiceItem struct {
   ProductID uint    `json:"product_id" binding:"required"`
   Quantity  int     `json:"quantity" binding:"required"`
   Price     float64 `json:"price" binding:"required"`
}

func (h *InvoiceHandler) Create(c *gin.Context) {
   var req createInvoiceRequest
   if err := c.ShouldBindJSON(&req); err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }

   var items []usecase.OrderItemInput
   for _, item := range req.Items {
       items = append(items, usecase.OrderItemInput{
           ProductID: item.ProductID,
           Quantity:  item.Quantity,
           Price:     item.Price,
       })
   }

   invoice, err := h.service.Create(c.Request.Context(), struct {
       UserID uint
       Items  []usecase.OrderItemInput
   }{
       UserID: req.UserID,
       Items:  items,
   })

   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusCreated, invoice)
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
   id, err := strconv.ParseUint(c.Param("id"), 10, 32)
   if err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
       return
   }

   invoice, err := h.service.GetByID(c.Request.Context(), uint(id))
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusOK, invoice)
}

type updateInvoiceRequest struct {
   Status string `json:"status" binding:"required"`
}

func (h *InvoiceHandler) Update(c *gin.Context) {
   id, err := strconv.ParseUint(c.Param("id"), 10, 32)
   if err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
       return
   }

   var req updateInvoiceRequest
   if err := c.ShouldBindJSON(&req); err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }

   invoice, err := h.service.Update(c.Request.Context(), uint(id), struct {
       Status string
   }{
       Status: req.Status,
   })

   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) Delete(c *gin.Context) {
   id, err := strconv.ParseUint(c.Param("id"), 10, 32)
   if err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
       return
   }

   if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

func (h *InvoiceHandler) List(c *gin.Context) {
   userID := c.Query("user_id")
   
   if userID != "" {
       id, err := strconv.ParseUint(userID, 10, 32)
       if err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
           return
       }
       
       invoices, err := h.service.GetByUserID(c.Request.Context(), uint(id))
       if err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
           return
       }
       
       c.JSON(http.StatusOK, invoices)
       return
   }

   invoices, err := h.service.List(c.Request.Context())
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusOK, invoices)
}

func (h *InvoiceHandler) GetUserInvoices(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
        return
    }

    invoices, err := h.service.GetUserInvoices(c.Request.Context(), userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, invoices)
}