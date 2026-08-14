package usecase

import (
	"context"
	"fmt"
	"strings"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/pkg/openfoodfacts"
	"t_dev_700/pkg/pagination"
)

type ProductService struct {
    repo      repository.ProductRepository
    offClient *openfoodfacts.Client
}

func NewProductService(repo repository.ProductRepository, offClient *openfoodfacts.Client) *ProductService {
    return &ProductService{
        repo:      repo,
        offClient: offClient,
    }
}

func (s *ProductService) Create(ctx context.Context, input struct {
    Name            string
    Price           float64
    Brand           string
    Picture         string
    Categories      []string
    NutritionalInfo string
    StockQuantity   int
    OpenFoodFactsID string
}) (*models.Product, error) {
    product := &models.Product{
        Name:            input.Name,
        Price:           input.Price,
        Brand:           input.Brand,
        Picture:         input.Picture,
        Categories:      input.Categories,
        NutritionalInfo: input.NutritionalInfo,
        StockQuantity:   input.StockQuantity,
        OpenFoodFactsID: input.OpenFoodFactsID,
    }

    if err := s.repo.Create(ctx, product); err != nil {
        return nil, err
    }

    return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, id uint, input struct {
    Name            string
    Price           float64
    Brand           string
    Picture         string
    Categories      []string
    NutritionalInfo string
    StockQuantity   int
    OpenFoodFactsID string
}) (*models.Product, error) {
    product, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    product.Name = input.Name
    product.Price = input.Price
    product.Brand = input.Brand
    product.Picture = input.Picture
    product.Categories = input.Categories
    product.NutritionalInfo = input.NutritionalInfo
    product.StockQuantity = input.StockQuantity
    product.OpenFoodFactsID = input.OpenFoodFactsID

    if err := s.repo.Update(ctx, product); err != nil {
        return nil, err
    }

    return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}

func (s *ProductService) List(ctx context.Context, page, pageSize int, filters *repository.ProductFilters) ([]models.Product, *pagination.Pagination, error) {
    p := pagination.NewPagination(page, pageSize)
    products, err := s.repo.List(ctx, p, filters)
    if err != nil {
        return nil, nil, err
    }
    return products, p, nil
}

func (s *ProductService) Search(ctx context.Context, query string, page, pageSize int) ([]models.Product, *pagination.Pagination, error) {
    p := pagination.NewPagination(page, pageSize)
    products, err := s.repo.Search(ctx, query, p)
    if err != nil {
        return nil, nil, err
    }
    return products, p, nil
}

func (s *ProductService) GetByCategory(ctx context.Context, category string, page, pageSize int) ([]models.Product, *pagination.Pagination, error) {
    p := pagination.NewPagination(page, pageSize)
    products, err := s.repo.GetByCategory(ctx, category, p)
    if err != nil {
        return nil, nil, err
    }
    return products, p, nil
}

func (s *ProductService) UpdateStock(ctx context.Context, id uint, quantity int) error {
    return s.repo.UpdateStock(ctx, id, quantity)
}

func (s *ProductService) CreateFromBarcode(ctx context.Context, barcode string, price float64, stockQuantity int) (*models.Product, error) {
    offProduct, err := s.offClient.GetProduct(barcode)
    if err != nil {
        return nil, err
    }

    categories := strings.Split(offProduct.Product.Categories, ",")
    for i := range categories {
        categories[i] = strings.TrimSpace(categories[i])
    }

    nutritionalInfo := fmt.Sprintf(
        "Proteins: %.2fg, Carbohydrates: %.2fg, Fat: %.2fg, Energy: %.2fkcal",
        offProduct.Product.Nutrients.Proteins100g,
        offProduct.Product.Nutrients.Carbohydrates100g,
        offProduct.Product.Nutrients.Fat100g,
        offProduct.Product.Nutrients.Energy100g,
    )

    product := &models.Product{
        Name:            offProduct.Product.ProductName,
        Price:           price,
        Brand:          offProduct.Product.Brands,
        Picture:        offProduct.Product.ImageURL,
        Categories:     categories,
        NutritionalInfo: nutritionalInfo,
        StockQuantity:   stockQuantity,
        OpenFoodFactsID: offProduct.Code,
    }

    if err := s.repo.Create(ctx, product); err != nil {
        return nil, err
    }

    return product, nil
}