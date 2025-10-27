package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware to recover from panics
func PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var panicMessage string
				if err, ok := r.(error); ok {
					panicMessage = err.Error()
				} else {
					panicMessage = fmt.Sprint(r)
				}
				respondError(c, http.StatusBadRequest, panicMessage)
			}
		}()

		c.Next()
	}
}
