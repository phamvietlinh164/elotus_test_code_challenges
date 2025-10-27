package middlewares

import (
	"hackathon/internal/common"
	"hackathon/internal/utils"
	"net/http"
	"strings"

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

		_, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if err.Error() == common.ErrTokenExpired {
				respondError(c, http.StatusUnauthorized, "Token is expired")
				return
			}
			respondError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, err := utils.GetClaimsFromJWT(tokenStr)
		if err != nil {
			respondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("claims", claims)
		c.Set("userID", claims.UserID)
		c.Set("isAdmin", claims.IsAdmin)

		if logger, exists := c.Get("logger"); exists {
			updatedLogger := logger.(zerolog.Logger).With().
				Uint("user_id", claims.UserID).
				Bool("is_admin", claims.IsAdmin).
				Logger()
			c.Set("logger", updatedLogger)
		}

		c.Next()
	}
}

func respondError(c *gin.Context, code int, msg string) {
	c.JSON(code, common.BaseApiResponse[any]{
		Success: false,
		Message: msg,
		Data:    nil,
	})
	c.Abort()
}
