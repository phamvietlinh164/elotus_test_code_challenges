package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request already has a Request-ID
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			// Generate a new UUID if none exists
			requestID = uuid.New().String()
		}

		// Add Request-ID to the context and response header
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDKey, requestID)

		c.Next()
	}
}
