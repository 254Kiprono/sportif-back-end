package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Configure allowed origins via CORS_ALLOW_ORIGINS
func CORSMiddleware() gin.HandlerFunc {
	allowed := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	allowAll := allowed == ""
	var allowList map[string]struct{}
	if !allowAll {
		allowList = make(map[string]struct{})
		for _, origin := range strings.Split(allowed, ",") {
			o := strings.TrimSpace(origin)
			if o != "" {
				allowList[o] = struct{}{}
			}
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			if _, ok := allowList[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Vary", "Origin")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
