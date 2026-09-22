package controller

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

// RequestID assigns a correlation id to every request so CLI errors can be
// matched with server logs without exposing internal stack traces.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" || len(id) > 64 {
			buf := make([]byte, 16)
			if _, err := rand.Read(buf); err != nil {
				id = "unknown"
			} else {
				id = hex.EncodeToString(buf)
			}
		}
		c.Set("requestID", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func BodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				ErrorFromCode(c, http.StatusInternalServerError, common.ErrorCodeInternal, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func ErrorResponse(c *gin.Context, status int, errorCode string, err error) {
	message := "request failed"
	if err != nil {
		message = err.Error()
	}
	writeFailCode(c, status, errorCode, message)
}

func ErrorFromCode(c *gin.Context, status int, errorCode, message string) {
	if errorCode == "" {
		errorCode = common.ErrorCodeInternal
	}
	writeFailCode(c, status, errorCode, message)
}
