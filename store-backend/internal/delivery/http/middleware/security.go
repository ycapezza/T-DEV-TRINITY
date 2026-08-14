package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func SetupSecurityMiddleware(r *gin.Engine) {
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:5173"}
    config.AllowCredentials = true
    config.AllowHeaders = []string{
        "Origin", 
        "Content-Type", 
        "Accept", 
        "Authorization",
        "X-CSRF-Token",
    }

    config.ExposeHeaders = []string{"X-CSRF-Token"}

    
    r.Use(cors.New(config))

    r.Use(secure.New(secure.Config{
        ContentSecurityPolicy: "default-src 'self'",
        ReferrerPolicy: "strict-origin-when-cross-origin",
        IsDevelopment: false,
        STSSeconds:            31536000,
        STSIncludeSubdomains: true,
        STSPreload:           true,
        FrameDeny:            true,
    }))
}