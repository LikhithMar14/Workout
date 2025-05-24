package auth


import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID int  `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}


func CustomClass(userID int,email string,issuer , audience string, expiration time.Duration ) *CustomClaims {
	now := time.Now()

	return &CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  []string{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}