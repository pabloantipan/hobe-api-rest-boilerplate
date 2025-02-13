package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/cloud"
	"github.com/pabloantipan/hobe-api-boilerplate/utils"
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
		payload := bodyDataCases(c)
		c.Next()
		end := time.Now()
		latency := end.Sub(start).Seconds()
		m.logger.LogRequest(c.Request.Method, payload, c.Request.URL.Path, c.Writer.Status(), latency)
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
		payload := bodyDataCases(c)
		c.Next()
		end := time.Now()
		latency := end.Sub(start).Seconds()

		if c.Writer.Status() >= 400 {
			if len(c.Errors) > 0 {
				if err := c.Errors.ByType(gin.ErrorTypePrivate).Last(); err != nil {
					m.logger.LogError(err, payload, c.Request.Method, c.Request.URL.Path, latency)
				} else {
					// Log error status without error object
					m.logger.LogError(fmt.Errorf("request failed with status %d", c.Writer.Status()),
						payload, c.Request.Method, c.Request.URL.Path, latency)
				}
			} else {
				// No specific error but status >= 400
				m.logger.LogError(fmt.Errorf("request failed with status %d", c.Writer.Status()),
					payload, c.Request.Method, c.Request.URL.Path, latency)
			}
		}

	}
}

func bodyDataCases(c *gin.Context) interface{} {
	if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
		body, err := utils.ExtractBody(c)
		if err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			// c.Abort()
			return map[string]interface{}{
				"headers": c.Request.Header,
				"body":    "Failed to parse request body" + fmt.Sprintf("%v", err),
				"query":   c.Request.URL.Query(),
			}
		}
		return map[string]interface{}{
			"headers": c.Request.Header,
			"body":    body,
			"query":   c.Request.URL.Query(),
		}
	}

	return map[string]interface{}{
		"headers": c.Request.Header,
		"body":    "No body data",
		"query":   c.Request.URL.Query(),
	}
}
