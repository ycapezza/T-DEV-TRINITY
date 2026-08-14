package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
   secretKey []byte
}

type Claims struct {
    UserID  uint   `json:"user_id"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"is_admin"`
    jwt.RegisteredClaims
}

func NewJWTAuth(secretKey string) *JWTAuth {
   return &JWTAuth{
       secretKey: []byte(secretKey),
   }
}

func (j *JWTAuth) GenerateToken(userID uint, email string, isAdmin bool) (string, error) {
    claims := &Claims{
        UserID:  userID,
        Email:   email,
        IsAdmin: isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
   claims := &Claims{}
   token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
       return j.secretKey, nil
   })

   if err != nil {
       return nil, err
   }

   if !token.Valid {
       return nil, errors.New("invalid token")
   }

   return claims, nil
}