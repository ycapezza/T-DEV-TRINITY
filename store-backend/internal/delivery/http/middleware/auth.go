package middleware

import (
	"net/http"
	"strings"
	"t_dev_700/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
    jwtAuth *auth.JWTAuth
}

func NewAuthMiddleware(jwtAuth *auth.JWTAuth) *AuthMiddleware {
    return &AuthMiddleware{
        jwtAuth: jwtAuth,
    }
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        claims, err := m.jwtAuth.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("is_admin", claims.IsAdmin)

        c.Next()
    }
}