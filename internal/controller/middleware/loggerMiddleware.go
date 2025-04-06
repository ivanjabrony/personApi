package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error("Controller error", slog.String("error", e.Error()))
			}
		}

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("incoming request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", statusCode),
			slog.String("client_ip", c.ClientIP()),
			slog.Duration("duration", duration),
		)
	}
}
