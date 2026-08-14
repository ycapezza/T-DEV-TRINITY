package usecase

import (
	"context"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) Create(ctx context.Context, input struct {
    FirstName   string
    LastName    string
    Email       string
    Password    string
    PhoneNumber string
    Address     string
    ZipCode     string
    City        string
    Country     string
}) (*models.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        FirstName:   input.FirstName,
        LastName:    input.LastName,
        Email:       input.Email,
        Password:    string(hashedPassword),
        PhoneNumber: input.PhoneNumber,
        Address:     input.Address,
        ZipCode:     input.ZipCode,
        City:        input.City,
        Country:     input.Country,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint, input struct {
    FirstName   string
    LastName    string
    PhoneNumber string
    Address     string
    ZipCode     string
    City        string
    Country     string
}) (*models.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    user.FirstName = input.FirstName
    user.LastName = input.LastName
    user.PhoneNumber = input.PhoneNumber
    user.Address = input.Address
    user.ZipCode = input.ZipCode
    user.City = input.City
    user.Country = input.Country

    if err := s.repo.Update(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
    return s.repo.List(ctx)
}