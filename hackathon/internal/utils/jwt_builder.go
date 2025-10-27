package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired = errors.New("token is expired")
)

// Init sets the secret key for JWT operations. Call once at app startup.
func Init(secret string) {
	secretKey = []byte(secret)
}

// GenerateToken creates HS256 token with iat/exp and custom fields.
func GenerateToken(userID string, isAdmin bool, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Subject:   userID,
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

// GetIsAdminFromJWT extracts is_admin
func GetIsAdminFromJWT(tokenStr string) (bool, error) {
	tok, err := ParseJWT(tokenStr)
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return false, ErrTokenExpired
		}
		return false, err
	}
	if claims, ok := tok.Claims.(*CustomClaims); ok && tok.Valid {
		return claims.IsAdmin, nil
	}
	return false, errors.New("invalid token claims")
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

var secretKey = []byte("your-secret-key")

type CustomClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GetUserIDFromJWT(tokenStr string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, errors.New("cannot parse claims")
	}

	id64, err := strconv.ParseUint(claims.UserID, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}

	return uint(id64), nil
}
