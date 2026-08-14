package usecase

import (
	"context"
	"errors"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"
	"t_dev_700/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo repository.UserRepository
    jwtAuth *auth.JWTAuth
}

func NewAuthService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtAuth:  jwtAuth,
    }
}

func (s *AuthService) Register(ctx context.Context, input struct {
    FirstName   string
    LastName    string
    Email       string
    Password    string
    PhoneNumber string
    Address     string
    ZipCode     string
    City        string
    Country     string
}) (*models.User, string, error) {
    if _, err := s.userRepo.GetByEmail(ctx, input.Email); err == nil {
        return nil, "", errors.New("email already exists")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, "", err
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

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, "", err
    }

    token, err := s.jwtAuth.GenerateToken(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

func (s *AuthService) LoginAdmin(ctx context.Context, email, password string) (*models.User, string, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    if !user.IsAdmin {
        return nil, "", errors.New("admin access required")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    token, err := s.jwtAuth.GenerateToken(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

func (s *AuthService) LoginMobile(ctx context.Context, email, password string) (*models.User, string, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    token, err := s.jwtAuth.GenerateToken(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}