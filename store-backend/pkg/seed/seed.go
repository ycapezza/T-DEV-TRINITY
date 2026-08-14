package seed

import (
	"fmt"
	"t_dev_700/internal/domain/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
    if err := db.Exec("DELETE FROM order_items").Error; err != nil {
        return err
    }
    if err := db.Exec("DELETE FROM invoices").Error; err != nil {
        return err
    }
    if err := db.Exec("DELETE FROM products").Error; err != nil {
        return err
    }
    if err := db.Exec("DELETE FROM users").Error; err != nil {
        return err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    adminUser := &models.User{
        FirstName:   "Admin",
        LastName:    "User",
        Email:      "admin@test.com",
        Password:    string(hashedPassword),
        PhoneNumber: "+1234567890",
        IsAdmin:    true,
    }
    if err := db.Create(adminUser).Error; err != nil {
        return err
    }

    regularUser := &models.User{
        FirstName:   "John",
        LastName:    "Doe",
        Email:      "john@test.com",
        Password:    string(hashedPassword),
        PhoneNumber: "+0987654321",
        IsAdmin:    false,
    }
    if err := db.Create(regularUser).Error; err != nil {
        return err
    }

    products := []struct {
        Name            string
        Price           float64
        Brand           string
        Categories      []string
        NutritionalInfo string
        StockQuantity   int
    }{
        {"Organic Milk", 3.99, "Happy Cow", []string{"dairy", "beverages", "organic"}, "Proteins: 3.4g, Carbohydrates: 4.8g, Fat: 3.6g", 50},
        {"Whole Grain Bread", 2.99, "Healthy Bakery", []string{"bakery", "whole grain", "bread"}, "Proteins: 7g, Carbohydrates: 30g, Fat: 1.5g", 30},
        {"Fresh Apples", 1.99, "Nature's Best", []string{"produce", "fruits", "fresh"}, "Carbohydrates: 14g, Fiber: 2.4g", 100},
        {"Chocolate Chip Cookies", 4.99, "Sweet Treats", []string{"snacks", "cookies", "sweet"}, "Proteins: 2g, Carbohydrates: 25g, Fat: 12g", 5},
        {"Orange Juice", 5.99, "Fresh Squeeze", []string{"beverages", "juice", "fruits"}, "Carbohydrates: 26g, Vitamin C: 100%", 25},
        {"Greek Yogurt", 4.49, "Mediterranean Delights", []string{"dairy", "yogurt", "protein"}, "Proteins: 15g, Carbohydrates: 6g, Fat: 0.5g", 40},
        {"Mixed Nuts", 7.99, "Nutty Goodness", []string{"snacks", "nuts", "protein"}, "Proteins: 6g, Fat: 14g, Fiber: 3g", 15},
        {"Quinoa", 6.99, "Ancient Grains", []string{"grains", "organic", "super foods"}, "Proteins: 8g, Carbohydrates: 39g, Fiber: 5g", 60},
    }

    var createdProducts []models.Product
    for _, p := range products {
        product := models.Product{
            Name:            p.Name,
            Price:           p.Price,
            Brand:           p.Brand,
            Categories:      p.Categories,
            NutritionalInfo: p.NutritionalInfo,
            StockQuantity:   p.StockQuantity,
        }
        
        if err := db.Create(&product).Error; err != nil {
            return err
        }
        createdProducts = append(createdProducts, product)
    }

    if len(createdProducts) == 0 {
        return fmt.Errorf("no products were created, cannot seed order items")
    }

    for i := 0; i < 10; i++ {
        invoice := models.Invoice{
            UserID: regularUser.ID,
            Status: "completed",
            CreatedAt: time.Now().AddDate(0, 0, -i), 
        }

        var total float64
        var orderItems []models.OrderItem

        for j := 0; j < (i%2)+2; j++ {
            productIndex := j % len(createdProducts)
            product := createdProducts[productIndex] 
            
            quantity := j + 1
            itemTotal := float64(quantity) * product.Price
            
            orderItems = append(orderItems, models.OrderItem{
                ProductID: product.ID,
                Quantity:  quantity,
                Price:    product.Price,
            })
            
            total += itemTotal
        }

        invoice.Total = total

        if err := db.Create(&invoice).Error; err != nil {
            return err
        }

        for _, item := range orderItems {
            item.InvoiceID = invoice.ID
            if err := db.Create(&item).Error; err != nil {
                return err
            }
        }
    }

    return nil
}
