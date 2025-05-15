package webauthn

import (
	"crypto/rand"
	"data-api/internal/handlers"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ensure WebAuthnHandler implements handlers.HandlerInterface
var _ handlers.HandlerInterface = (*WebAuthnHandler)(nil)

// In-memory store for demonstration (replace with DB/Redis in production)
var (
	userCredentials = make(map[string]interface{})
)

type WebAuthnHandler struct{}

func NewWebAuthnHandler() *WebAuthnHandler {
	return &WebAuthnHandler{}
}

func (h *WebAuthnHandler) SetupRoutes(rg *gin.RouterGroup) {
	webauthn := rg.Group("/webauthn")
	webauthn.POST("/register/options", h.RegisterOptions)
	webauthn.POST("/register/verify", h.RegisterVerify)
	webauthn.POST("/login/options", h.LoginOptions)
	webauthn.POST("/login/verify", h.LoginVerify)
}

func randomB64(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (h *WebAuthnHandler) RegisterOptions(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	// Return a valid WebAuthn PublicKeyCredentialCreationOptions JSON
	options := gin.H{
		"challenge": randomB64(32),
		"rp":        gin.H{"name": "Data API", "id": "localhost"},
		"user": gin.H{
			"id":          base64.RawURLEncoding.EncodeToString([]byte(req.Username)),
			"name":        req.Username,
			"displayName": req.Username,
		},
		"pubKeyCredParams": []gin.H{{"type": "public-key", "alg": -7}},
		"timeout":          60000,
		"attestation":      "none",
	}
	c.JSON(http.StatusOK, options)
}

func (h *WebAuthnHandler) RegisterVerify(c *gin.Context) {
	var req struct {
		Username    string      `json:"username"`
		Attestation interface{} `json:"attestation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// Store attestation (replace with real verification)
	userCredentials[req.Username] = req.Attestation
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *WebAuthnHandler) LoginOptions(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	// Return dummy options (replace with real WebAuthn options)
	options := gin.H{
		"challenge":        "dummy-challenge",
		"allowCredentials": []gin.H{{"id": req.Username, "type": "public-key"}},
	}
	c.JSON(http.StatusOK, options)
}

func (h *WebAuthnHandler) LoginVerify(c *gin.Context) {
	var req struct {
		Username  string      `json:"username"`
		Assertion interface{} `json:"assertion"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// Dummy check (replace with real assertion verification)
	if _, ok := userCredentials[req.Username]; !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": fmt.Sprintf("User %s authenticated", req.Username)})
}
