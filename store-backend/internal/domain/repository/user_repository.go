package repository

import (
	"context"
	"t_dev_700/internal/domain/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context) ([]models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
}