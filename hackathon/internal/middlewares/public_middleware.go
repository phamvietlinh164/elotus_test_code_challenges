package middlewares

// import (
// 	"hackathon/internal/utils"
// 	"net/http"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/zerolog"
// )

// func Public() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")

// 		if authHeader == "" {
// 			c.Next()
// 			return
// 		}

// 		tokenString := authHeader
// 		tokenNoBear := utils.GetTokenFromBearer(tokenString)
// 		token, err := utils.ParseJWT(tokenNoBear)
// 		if err != nil || !token.Valid {
// 			//if err.Error() == common.ErrTokenExpired {
// 			//	c.JSON(http.StatusUnauthorized, common.BaseApiResponse[any]{
// 			//		HttpRequestStatus: http.StatusUnauthorized,
// 			//		Success:           false,
// 			//		Message:           "Token is expired",
// 			//		Data:              nil,
// 			//	})
// 			//	c.Abort()
// 			//	return
// 			//}
// 			//
// 			//c.JSON(http.StatusBadRequest, common.BaseApiResponse[any]{
// 			//	HttpRequestStatus: http.StatusBadRequest,
// 			//	Success:           false,
// 			//	Message:           "Invalid token",
// 			//	Data:              nil,
// 			//})
// 			//c.Abort()
// 			c.Next()
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			c.Set("claims", claims)
// 		}

// 		userId, err := utils.GetUserIDFromJWT(tokenNoBear)
// 		if err != nil {
// 			respondError(c, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		isAdmin, err := utils.GetIsAdminFromJWT(tokenNoBear)
// 		if err != nil {
// 			respondError(c, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		c.Set("userID", userId)
// 		c.Set("isAdmin", isAdmin)

// 		if logger, exists := c.Get("logger"); exists {
// 			updatedLogger := logger.(zerolog.Logger).With().Uint("user_id", userId).Bool("is_admin", isAdmin).Logger()
// 			c.Set("logger", updatedLogger)
// 		}
// 		c.Next()
// 	}
// }
