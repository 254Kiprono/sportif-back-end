package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"webuye-sportif/app/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// AuthMiddleware validates the JWT and, if Redis is available, verifies the session
// is still live (whitelisting approach). rdb can be nil — session check is skipped.
func AuthMiddleware(cfg *config.Config, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Validate JTI (session ID) against Redis whitelist
		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing session identifier"})
			c.Abort()
			return
		}

		if rdb != nil {
			ctx := context.Background()
			_, err = rdb.Get(ctx, fmt.Sprintf("session:%s", jti)).Result()
			if err != nil {
				// Key missing = session was logged out or expired in Redis
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or terminated. Please login again."})
				c.Abort()
				return
			}
		}

		// Attach claims to context for downstream handlers
		c.Set("jti", jti)
		c.Set("user_id", claims["user_id"])
		c.Set("role_id", claims["role_id"])
		c.Set("role_name", claims["role_name"])
		if perms, ok := claims["permissions"].([]interface{}); ok {
			var permissions []string
			for _, p := range perms {
				if s, ok := p.(string); ok {
					permissions = append(permissions, s)
				}
			}
			c.Set("permissions", permissions)
		}
		c.Next()
	}
}

func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role_name")
		if currentUserRole == "admin" {
			c.Next()
			return
		}

		if currentUserRole != roleName {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions: role " + roleName + " required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAnyRole(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role_name")
		if currentUserRole == "admin" {
			c.Next()
			return
		}

		for _, r := range roles {
			if currentUserRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions: one of the required roles not found"})
		c.Abort()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role_name")
		if currentUserRole == "admin" {
			c.Next()
			return
		}

		perms, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		permissions := perms.([]string)
		for _, p := range permissions {
			if p == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions: " + permission + " required"})
		c.Abort()
	}
}
