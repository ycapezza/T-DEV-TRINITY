package openfoodfacts

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
}

type Product struct {
    Code       string `json:"code"`
    Status     int    `json:"status"`
    Product    struct {
        ProductName    string `json:"product_name"`
        Brands        string `json:"brands"`
        ImageURL      string `json:"image_url"`
        Categories    string `json:"categories"`
        Nutrients     struct {
            Proteins100g   float64 `json:"proteins_100g"`
            Carbohydrates100g float64 `json:"carbohydrates_100g"`
            Fat100g         float64 `json:"fat_100g"`
            Energy100g     float64 `json:"energy-kcal_100g"`
        } `json:"nutriments"`
    } `json:"product"`
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{},
        baseURL:    "https://world.openfoodfacts.org/api/v0",
    }
}

func (c *Client) GetProduct(barcode string) (*Product, error) {
    url := fmt.Sprintf("%s/product/%s.json", c.baseURL, barcode)
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var product Product
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, err
    }

    if product.Status != 1 {
        return nil, fmt.Errorf("product not found")
    }

    return &product, nil
}