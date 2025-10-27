package middlewares

import (
	"hackathon/internal/common"
	"hackathon/internal/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondError(c, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		tokenStr := utils.GetTokenFromBearer(authHeader)
		if strings.TrimSpace(tokenStr) == "" {
			respondError(c, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		token, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if err.Error() == common.ErrTokenExpired {
				respondError(c, http.StatusUnauthorized, "Token is expired")
				return
			}
			respondError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respondError(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID, err := utils.GetUserIDFromJWT(tokenStr)
		if err != nil {
			respondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		isAdmin, err := utils.GetIsAdminFromJWT(tokenStr)
		if err != nil {
			respondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("claims", claims)
		c.Set("userID", userID)
		c.Set("isAdmin", isAdmin)

		if logger, exists := c.Get("logger"); exists {
			updatedLogger := logger.(zerolog.Logger).With().
				Uint("user_id", userID).
				Bool("is_admin", isAdmin).
				Logger()
			c.Set("logger", updatedLogger)
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("isAdmin")
		if !exists {
			respondError(c, http.StatusUnauthorized, "Missing admin flag in context")
			return
		}

		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			respondError(c, http.StatusForbidden, "Admin access required")
			return
		}

		c.Next()
	}
}

func respondError(c *gin.Context, code int, msg string) {
	c.JSON(code, common.BaseApiResponse[any]{
		HttpRequestStatus: code,
		Success:           false,
		Message:           msg,
		Data:              nil,
	})
	c.Abort()
}
