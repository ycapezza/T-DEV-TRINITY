package handler

import (
	"net/http"
	"t_dev_700/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    service *usecase.AuthService
}

func NewAuthHandler(service *usecase.AuthService) *AuthHandler {
    return &AuthHandler{
        service: service,
    }
}

type registerRequest struct {
    FirstName   string `json:"first_name" binding:"required"`
    LastName    string `json:"last_name" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    PhoneNumber string `json:"phone_number"`
    Address     string `json:"address"`
    ZipCode     string `json:"zip_code"`
    City        string `json:"city"`
    Country     string `json:"country"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.service.Register(c.Request.Context(), struct {
        FirstName   string
        LastName    string
        Email       string
        Password    string
        PhoneNumber string
        Address     string
        ZipCode     string
        City        string
        Country     string
    }{
        FirstName:   req.FirstName,
        LastName:    req.LastName,
        Email:       req.Email,
        Password:    req.Password,
        PhoneNumber: req.PhoneNumber,
        Address:     req.Address,
        ZipCode:     req.ZipCode,
        City:        req.City,
        Country:     req.Country,
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "user":  user,
        "token": token,
    })
}

type loginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) LoginAdmin(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.service.LoginAdmin(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user":  user,
        "token": token,
    })
}

func (h *AuthHandler) LoginMobile(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.service.LoginMobile(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user":  user,
        "token": token,
    })
}