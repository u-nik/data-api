package middleware

import (
	"context"
	"data-api/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	verifier *oidc.IDTokenVerifier
	clientID = utils.GetRequiredEnv("OIDC_CLIENT_ID")
)

type TokenClaims struct {
	Subject   string   `json:"sub"`
	ClientID  string   `json:"client_id"`
	Scope     string   `json:"scope"`
	Scp       []string `json:"scp"`
	Roles     []string `json:"roles,omitempty"` // optional
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
}

func init() {
	log.Printf("Initializing OIDC provider with issuer: %s", utils.GetRequiredEnv("OIDC_ISSUER"))
	// Initialize the OIDC provider
	provider, err := oidc.NewProvider(context.Background(), utils.GetRequiredEnv("OIDC_ISSUER"))
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	// Create an ID token verifier
	verifier = provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})
}

func Auth(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		baseLogger.Sugar().Infof("Authorize Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(c.Request.Context(), rawToken)
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

func GetTokenClaims(c *gin.Context) (TokenClaims, bool) {
	val, exists := c.Get("claims")
	if !exists {
		return TokenClaims{}, false
	}

	claims, ok := val.(TokenClaims)
	return claims, ok
}

func GetScopes(claims TokenClaims) []string {
	if len(claims.Scp) > 0 {
		return claims.Scp
	}

	if claims.Scope != "" {
		return strings.Split(claims.Scope, " ")
	}

	return nil
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
