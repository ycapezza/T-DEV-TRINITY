package postgres

import (
	"context"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/pkg/pagination"

	"gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{
        db: db,
    }
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).First(&product, id).Error
    return &product, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, p *pagination.Pagination, filters *repository.ProductFilters) ([]models.Product, error) {
    var products []models.Product
    query := r.db.WithContext(ctx)

    if filters != nil {
        if filters.Name != "" {
            query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
        }
        if filters.Category != "" {
            query = query.Where("? = ANY(categories)", filters.Category)
        }
    }

    var total int64
    if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
        return nil, err
    }
    p.Total = total

    err := query.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&products).Error
    if err != nil {
        return nil, err
    }

    return products, nil
}

func (r *productRepository) Search(ctx context.Context, query string, p *pagination.Pagination) ([]models.Product, error) {
    var products []models.Product
    db := r.db.WithContext(ctx)

    searchQuery := db.Where("name ILIKE ?", "%"+query+"%").
        Or("brand ILIKE ?", "%"+query+"%").
        Or("? = ANY(categories)", query)

    var total int64
    if err := searchQuery.Model(&models.Product{}).Count(&total).Error; err != nil {
        return nil, err
    }
    p.Total = total

    err := searchQuery.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&products).Error
    if err != nil {
        return nil, err
    }

    return products, nil
}

func (r *productRepository) GetByCategory(ctx context.Context, category string, p *pagination.Pagination) ([]models.Product, error) {
    var products []models.Product
    query := r.db.WithContext(ctx).Where("? = ANY(categories)", category)

    var total int64
    if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
        return nil, err
    }
    p.Total = total

    err := query.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&products).Error
    if err != nil {
        return nil, err
    }

    return products, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
    return r.db.WithContext(ctx).Model(&models.Product{}).
        Where("id = ?", id).
        Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}