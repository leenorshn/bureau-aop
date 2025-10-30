package auth

import (
	"errors"
	"time"

	"bureau/internal/config"
	"bureau/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTService struct {
	config *config.Config
	logger *zap.Logger
}

func NewJWTService(cfg *config.Config, logger *zap.Logger) *JWTService {
	return &JWTService{
		config: cfg,
		logger: logger,
	}
}

type Claims struct {
	AdminID string `json:"admin_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new access token
func (j *JWTService) GenerateAccessToken(admin *models.Admin) (string, error) {
	claims := &Claims{
		AdminID: admin.ID.Hex(),
		Email:   admin.Email,
		Role:    admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.JWTAccessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bureau-mlm",
			Subject:   admin.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.JWTSecret))
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTService) GenerateRefreshToken(admin *models.Admin) (string, error) {
	claims := &Claims{
		AdminID: admin.ID.Hex(),
		Email:   admin.Email,
		Role:    admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.JWTRefreshExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bureau-mlm",
			Subject:   admin.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.JWTRefreshSecret))
}

// ValidateAccessToken validates an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.JWTRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// ExtractClaimsFromToken extracts claims from a token without validation
func (j *JWTService) ExtractClaimsFromToken(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token format")
}
