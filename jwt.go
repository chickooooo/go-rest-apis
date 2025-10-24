package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("JSONWebToken")

const accessTokenTTL = 15 * time.Minute
const refreshTokenTTL = 24 * time.Hour

// Claims defines the structure for JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// generateToken creates a signed JWT for the given userID and TTL
func generateToken(userID int, ttl time.Duration) (string, error) {
	// Create token claims
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	// Create, sign and return the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func NewAccessToken(userID int) (string, error) {
	return generateToken(userID, accessTokenTTL)
}

func NewRefreshToken(userID int) (string, error) {
	return generateToken(userID, refreshTokenTTL)
}

// VerifyToken parses and validates a JWT, returning the user ID if valid.
func VerifyToken(tokenStr string) (int, error) {
	// Parse the token with a key lookup function
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// Ensure the signing method is HMAC and expected
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenInvalidClaims
			}
			return jwtSecret, nil
		})
	if err != nil {
		return 0, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, jwt.ErrTokenInvalidClaims
}
