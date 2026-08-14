package bootstrap

import (
	"context"
	"t_dev_700/internal/domain/models"
	"t_dev_700/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

func CreateInitialAdmin(userRepo repository.UserRepository) error {
    ctx := context.Background()

    admins, err := userRepo.List(ctx)
    if err != nil {
        return err
    }

    for _, user := range admins {
        if user.IsAdmin {
            return nil
        }
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    initialAdmin := &models.User{
        FirstName:   "Admin",
        LastName:    "User",
        Email:       "admin@test.com",
        Password:    string(hashedPassword),
        PhoneNumber: "0000000000",
        IsAdmin:     true,
    }

    return userRepo.Create(ctx, initialAdmin)
}