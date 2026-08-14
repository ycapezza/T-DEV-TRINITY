package repository

import (
	"context"
	"t_dev_700/internal/domain/models"
	"t_dev_700/pkg/pagination"
)

type ProductFilters struct {
    Name     string
    Category string
}

type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) error
    GetByID(ctx context.Context, id uint) (*models.Product, error)
    Update(ctx context.Context, product *models.Product) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, p *pagination.Pagination, filters *ProductFilters) ([]models.Product, error)
    UpdateStock(ctx context.Context, id uint, quantity int) error
    GetByCategory(ctx context.Context, category string, p *pagination.Pagination) ([]models.Product, error)
    Search(ctx context.Context, query string, p *pagination.Pagination) ([]models.Product, error)
}