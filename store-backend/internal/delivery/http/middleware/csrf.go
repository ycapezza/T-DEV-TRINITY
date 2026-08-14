package middleware

import (
	"net/http"
	"t_dev_700/pkg/csrf"

	"github.com/gin-gonic/gin"
)

type CSRFMiddleware struct {
    tokenStore *csrf.TokenStore
}

func NewCSRFMiddleware() *CSRFMiddleware {
    return &CSRFMiddleware{
        tokenStore: csrf.NewTokenStore(),
    }
}

func (m *CSRFMiddleware) GenerateToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := m.tokenStore.Generate()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
            c.Abort()
            return
        }

        c.Header("X-CSRF-Token", token)
        c.Next()
    }
}

func (m *CSRFMiddleware) ValidateToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "GET" {
            c.Next()
            return
        }

        token := c.GetHeader("X-CSRF-Token")
        if token == "" {
            c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
            c.Abort()
            return
        }

        if !m.tokenStore.Verify(token) {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired CSRF token"})
            c.Abort()
            return
        }

        c.Next()
    }
}