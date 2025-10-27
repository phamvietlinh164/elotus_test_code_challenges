package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = "unknown"
		}

		logger := log.With().
			Timestamp().
			Str("request_id", requestID).
			Str("method", method).
			Str("path", path).
			Str("client_ip", clientIP).
			Str("user_agent", userAgent).
			Logger()

		c.Set("logger", logger)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		event := logger.Info()
		if status >= 400 {
			event = logger.Error()
		}

		event.
			Int("status", status).
			Dur("duration", duration).
			Msg("request completed")
	}
}
