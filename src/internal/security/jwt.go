package security

import (
	"errors"
	"log/slog"
	"time"

	configuration "personal-portfolio-main-back/src/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TODO - Everything related with JWT is pretty useless right now - I'll come back to it later
type JWTClaims struct {
	jwt.RegisteredClaims
	ClientType string `json:"client_type"`
}

type JWTManager struct {
	signingKey []byte
	config     configuration.Config
}

func NewJWTManager(config configuration.Config) *JWTManager {
	return &JWTManager{
		signingKey: []byte(config.JWTSigningKey),
		config:     config,
	}
}

func (j *JWTManager) GenerateToken() (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.config.JWTExpirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "frontend-client",
		},
		ClientType: "frontend",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.signingKey)
	if err != nil {
		slog.Error("Token generation failed")
		return "", errors.New("token generation failed")
	}

	return signedToken, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("Invalid JWT signing method detected")
			return nil, errors.New("invalid token")
		}
		return j.signingKey, nil
	})

	if err != nil {
		slog.Warn("JWT token validation failed")
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	slog.Warn("JWT token validation failed", "reason", "invalid_claims")
	return nil, errors.New("invalid token")
}
