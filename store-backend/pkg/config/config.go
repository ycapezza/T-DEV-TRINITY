package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    Port string

    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string

    JWTSecretKey string

    OpenFoodFactsAPIKey string
    PaypalClientID      string
    PaypalClientSecret  string
    PaypalWebhookID    string

}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return &Config{
        Port: getEnv("PORT", "8080"),

        DBHost:     getEnv("DB_HOST", "localhost"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "t_dev_700"),
        DBPort:     getEnv("DB_PORT", "5432"),

        JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),

        OpenFoodFactsAPIKey: getEnv("OPENFOODFACTS_API_KEY", ""),
        PaypalClientID:      getEnv("PAYPAL_CLIENT_ID", ""),
        PaypalClientSecret:  getEnv("PAYPAL_CLIENT_SECRET", ""),
        PaypalWebhookID:    getEnv("PAYPAL_WEBHOOK_ID", ""),
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}