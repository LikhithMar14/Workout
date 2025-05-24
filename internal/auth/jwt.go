package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator is a helper for generating and validating JWTs.
type JWTAuthenticator struct {
	secret string // HMAC secret used to sign/verify tokens
	aud    string // Expected audience ("who the token is for")
	iss    string // Expected issuer ("who issued the token")
}

// NewJWTAuthenticator creates a new instance with the given secret, audience, and issuer.
// Proper order: secret -> aud -> iss
func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

// GenerateToken creates a JWT token using the provided claims.
// The claims must include exp, iss, aud, sub etc., usually using jwt.RegisteredClaims.
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	// Create a new token with signing method HS256 and the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err // return error if signing fails
	}

	return tokenString, nil // return signed JWT string
}

// ValidateToken verifies the token's signature and standard claims like exp, aud, iss, alg.
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token,
		// Key function: provides the secret key after checking algorithm
		func(t *jwt.Token) (any, error) {
			// Validate that the signing method is HMAC (e.g. HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			}

			// Return the secret key to be used for signature validation
			return []byte(a.secret), nil
		},

		// Additional validation options for claims

		// Require "exp" (expiration time) claim to be present and not expired
		jwt.WithExpirationRequired(),

		// Validate the "aud" (audience) claim — ensures the token was intended for this server/service
		jwt.WithAudience(a.aud),

		// Validate the "iss" (issuer) claim — ensures the token was issued by a trusted authority
		jwt.WithIssuer(a.iss),

		// Accept only this specific signing method to avoid downgrade or substitution attacks
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
