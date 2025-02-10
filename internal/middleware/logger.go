package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/cloud"
)

type RequestLoggerMiddleware struct {
	logger *cloud.CloudLogger
}

func NewRequestLoggerMiddleware(logger *cloud.CloudLogger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		logger: logger,
	}
}

func (m *RequestLoggerMiddleware) HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start).Seconds()
		m.logger.LogRequest(c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

type ResponseLoggerMiddleware struct {
	logger *cloud.CloudLogger
}

func NewResponseLoggerMiddleware(logger *cloud.CloudLogger) *ResponseLoggerMiddleware {
	return &ResponseLoggerMiddleware{
		logger: logger,
	}
}

func (m *ResponseLoggerMiddleware) HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start).Seconds()

		if c.Writer.Status() >= 400 {
			m.logger.LogError(c.Errors.ByType(gin.ErrorTypePrivate).Last(), c.Request.Method, c.Request.URL.Path, latency)
		}

	}
}
