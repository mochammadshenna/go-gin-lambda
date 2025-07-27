package authentication

import (
	"ai-service/internal/util/exceptioncode"
	"ai-service/internal/util/logger"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret string

// Init initializes the JWT secret
func Init(secret string) {
	jwtSecret = secret
}

// GenerateToken generates a JWT token
func GenerateToken(userID, username, email, role string) (string, error) {
	claims := JWTClaim{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logger.Error(context.Background(), "token validation failed", err)
		return nil, exceptioncode.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, exceptioncode.ErrTokenInvalid
}

// ExtractToken extracts token from Authorization header
func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", exceptioncode.ErrTokenInvalid
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", exceptioncode.ErrTokenInvalid
	}

	return parts[1], nil
}

// ExtractClaim extracts JWT claims from request
func ExtractClaim(authHeader string) (*JWTClaim, error) {
	tokenString, err := ExtractToken(authHeader)
	if err != nil {
		return nil, err
	}

	return ValidateToken(tokenString)
}

// IsTokenExpired checks if token is expired
func IsTokenExpired(tokenString string) bool {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(claims.ExpiresAt.Time)
}

// RefreshToken refreshes a JWT token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Check if token is close to expiration (within 1 hour)
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", fmt.Errorf("token is not close to expiration")
	}

	// Generate new token with same claims but new expiration
	return GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role)
}
