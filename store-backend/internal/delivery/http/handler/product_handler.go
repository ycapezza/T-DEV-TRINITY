package handler

import (
	"net/http"
	"strconv"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
    service *usecase.ProductService
}

func NewProductHandler(service *usecase.ProductService) *ProductHandler {
    return &ProductHandler{
        service: service,
    }
}

type createProductRequest struct {
    Name            string    `json:"name" binding:"required"`
    Price           float64   `json:"price" binding:"required"`
    Brand           string    `json:"brand" binding:"required"`
    Picture         string    `json:"picture"`
    Categories      []string  `json:"categories" binding:"required"`
    NutritionalInfo string    `json:"nutritional_info"`
    StockQuantity   int       `json:"stock_quantity" binding:"required"`
    OpenFoodFactsID string    `json:"open_food_facts_id"`
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req createProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.service.Create(c.Request.Context(), struct {
        Name            string
        Price           float64
        Brand           string
        Picture         string
        Categories      []string
        NutritionalInfo string
        StockQuantity   int
        OpenFoodFactsID string
    }{
        Name:            req.Name,
        Price:           req.Price,
        Brand:           req.Brand,
        Picture:         req.Picture,
        Categories:      req.Categories,
        NutritionalInfo: req.NutritionalInfo,
        StockQuantity:   req.StockQuantity,
        OpenFoodFactsID: req.OpenFoodFactsID,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    product, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, product)
}

type updateProductRequest struct {
    Name            string    `json:"name" binding:"required"`
    Price           float64   `json:"price" binding:"required"`
    Brand           string    `json:"brand" binding:"required"`
    Picture         string    `json:"picture"`
    Categories      []string  `json:"categories" binding:"required"`
    NutritionalInfo string    `json:"nutritional_info"`
    StockQuantity   int       `json:"stock_quantity" binding:"required"`
    OpenFoodFactsID string    `json:"open_food_facts_id"`
}

func (h *ProductHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var req updateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.service.Update(c.Request.Context(), uint(id), struct {
        Name            string
        Price           float64
        Brand           string
        Picture         string
        Categories      []string
        NutritionalInfo string
        StockQuantity   int
        OpenFoodFactsID string
    }{
        Name:            req.Name,
        Price:           req.Price,
        Brand:           req.Brand,
        Picture:         req.Picture,
        Categories:      req.Categories,
        NutritionalInfo: req.NutritionalInfo,
        StockQuantity:   req.StockQuantity,
        OpenFoodFactsID: req.OpenFoodFactsID,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    
    filters := &repository.ProductFilters{
        Name:     c.Query("name"),
        Category: c.Query("category"),
    }

    products, p, err := h.service.List(c.Request.Context(), page, pageSize, filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "pagination": p,
    })
}

func (h *ProductHandler) Search(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    products, p, err := h.service.Search(c.Request.Context(), query, page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "pagination": p,
    })
}

type createProductFromBarcodeRequest struct {
    Barcode       string  `json:"barcode" binding:"required"`
    Price         float64 `json:"price" binding:"required"`
    StockQuantity int     `json:"stock_quantity" binding:"required"`
}

func (h *ProductHandler) CreateFromBarcode(c *gin.Context) {
    var req createProductFromBarcodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.service.CreateFromBarcode(
        c.Request.Context(),
        req.Barcode,
        req.Price,
        req.StockQuantity,
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}