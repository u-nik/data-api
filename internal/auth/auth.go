package auth

import (
	"context"
	"data-api/internal/utils"
	"log"
	"strings"

	goOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Verifier   *goOidc.IDTokenVerifier
	clientID   = utils.GetEnv("OIDC_CLIENT_ID", "data-api")
	oidcIssuer = utils.GetEnv("OIDC_ISSUER", "") // OIDC issuer URL
)

type TokenClaims struct {
	Subject   string   `json:"sub"`
	ClientID  string   `json:"client_id"`
	TanentID  string   `json:"tenant_id"`
	Scope     string   `json:"scope"`
	Scp       []string `json:"scp"`
	Roles     []string `json:"roles,omitempty"` // optional
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
}

func Initialize() {
	if oidcIssuer != "" {
		zap.L().Sugar().Infof("Initializing OIDC provider with issuer: %s", oidcIssuer)
		// Initialize the OIDC provider
		provider, err := goOidc.NewProvider(context.Background(), oidcIssuer)
		if err != nil {
			log.Fatalf("Failed to create OIDC provider: %v", err)
		}

		// Create an ID token verifier
		Verifier = provider.Verifier(&goOidc.Config{
			ClientID: clientID,
		})
	} else {
		zap.L().Sugar().Infoln("OIDC issuer not set, skipping OIDC initialization")
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
