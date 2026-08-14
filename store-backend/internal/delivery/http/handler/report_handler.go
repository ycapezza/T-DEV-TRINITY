package handler

import (
	"net/http"
	"strconv"
	"t_dev_700/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
    service *usecase.ReportService
}

func NewReportHandler(service *usecase.ReportService) *ReportHandler {
    return &ReportHandler{
        service: service,
    }
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")

    start, err := time.Parse("2006-01-02", startDate)
    if err != nil {
        start = time.Now().AddDate(0, -1, 0)
    }

    end, err := time.Parse("2006-01-02", endDate)
    if err != nil {
        end = time.Now()
    }

    report, err := h.service.GenerateSalesReport(c.Request.Context(), start, end)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetTopProducts(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    reports, err := h.service.GetTopProducts(c.Request.Context(), limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) GetCategoryPerformance(c *gin.Context) {
    reports, err := h.service.GetCategoryPerformance(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) GetStockAlerts(c *gin.Context) {
    minStock, _ := strconv.Atoi(c.DefaultQuery("minimum_stock", "10"))
    
    alerts, err := h.service.GetStockAlerts(c.Request.Context(), minStock)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, alerts)
}

func (h *ReportHandler) GetSalesEvolution(c *gin.Context) {
    period := c.DefaultQuery("period", "daily")
    
    evolution, err := h.service.GetSalesEvolution(c.Request.Context(), period)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, evolution)
}