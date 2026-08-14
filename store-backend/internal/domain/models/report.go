package models

type SalesReport struct {
    TotalRevenue     float64 `json:"total_revenue"`
    TotalOrders      int     `json:"total_orders"`
    AverageOrderSize float64 `json:"average_order_size"`
}

type ProductSalesReport struct {
    ProductID      uint    `json:"product_id"`
    ProductName    string  `json:"product_name"`
    QuantitySold   int     `json:"quantity_sold"`
    TotalRevenue   float64 `json:"total_revenue"`
}

type CategoryReport struct {
    Category     string  `json:"category"`
    TotalSales   float64 `json:"total_sales"`
    OrderCount   int     `json:"order_count"`
}

type StockAlert struct {
    ProductID      uint    `json:"product_id"`
    ProductName    string  `json:"product_name"`
    CurrentStock   int     `json:"current_stock"`
    MinimumStock   int     `json:"minimum_stock"`
}

type SalesEvolution struct {
    Period    string  `json:"period"`
    Revenue   float64 `json:"revenue"`
    OrderCount int    `json:"order_count"`
}