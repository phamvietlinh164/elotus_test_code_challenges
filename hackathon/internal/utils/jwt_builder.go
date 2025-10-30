package utils

import (
	"errors"
	"strings"
	"time"

	"hackathon/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	secretKey       []byte
)

func Init() {
	secretKey = []byte(config.Cfg.Jwt.Secret)
}

// GenerateToken creates HS256 token with iat/exp and custom fields.
func GenerateToken(userID uint, isAdmin bool, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseJWT parses token string and returns *jwt.Token or error.
// If token is expired, returns ErrTokenExpired
func ParseJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// ensure signing method is HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// optional: apply manual leeway (e.g. 5s tolerance)
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		now := time.Now().UTC()
		leeway := 5 * time.Second

		// Check expiration with leeway
		if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time.Add(leeway)) {
			return nil, ErrTokenExpired
		}

		// Check issued-at validity (not in future)
		if claims.IssuedAt != nil && now.Before(claims.IssuedAt.Time.Add(-leeway)) {
			return nil, errors.New("token used before issued")
		}
	}

	return token, nil
}

func GetTokenFromBearer(header string) string {
	// header may be "Bearer <token>" or just "<token>"
	parts := strings.Fields(header)
	if len(parts) == 0 {
		return ""
	}
	if strings.ToLower(parts[0]) == "bearer" && len(parts) >= 2 {
		return parts[1]
	}
	// fallback
	return parts[0]
}

type CustomClaims struct {
	UserID  uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func GetClaimsFromJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}
	return claims, nil
}
