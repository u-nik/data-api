package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if oidcIssuer == "" {
			c.Next()
			return
		}

		zap.L().Sugar().Infof("Authorize Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := Verifier.Verify(c.Request.Context(), rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var claims TokenClaims
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RequireScope(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if oidcIssuer == "" {
			c.Next()
			return
		}

		claims, exists := GetTokenClaims(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no claims"})
			return
		}

		for _, scope := range GetScopes(claims) {
			if scope == required {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing required scope " + required})
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if oidcIssuer == "" {
			c.Next()
			return
		}

		claims, ok := GetTokenClaims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing claims"})
			return
		}

		for _, r := range claims.Roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}
