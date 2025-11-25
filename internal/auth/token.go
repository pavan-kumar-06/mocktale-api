package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
}

var secretKey = []byte("your-secret-key-change-in-production")

// GenerateAndSetToken creates new token and sets HTTP-only cookie
func GenerateAndSetToken(w http.ResponseWriter) (string, error) {
	userID := uuid.New().String()
	now := time.Now()
	
	payload := TokenPayload{
		UserID:    userID,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(365 * 24 * time.Hour).Unix(),
	}
	
	payloadJSON, _ := json.Marshal(payload)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(payloadJSON)
	signature := mac.Sum(nil)
	
	tokenData := base64.URLEncoding.EncodeToString(payloadJSON) + "." + base64.URLEncoding.EncodeToString(signature)
	
	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    tokenData,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	
	return tokenData, nil
}

// ValidateToken validates the token from cookie
func ValidateToken(r *http.Request) (*TokenPayload, bool) {
	cookie, err := r.Cookie("user_token")
	if err != nil {
		return nil, false
	}
	
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return nil, false
	}
	
	payloadJSON, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, false
	}
	
	// Verify signature
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(payloadJSON)
	expectedSignature := mac.Sum(nil)
	
	actualSignature, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, false
	}
	
	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, false
	}
	
	// Parse payload
	var payload TokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, false
	}
	
	// Check expiration
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, false
	}
	
	return &payload, true
}

// GetOrCreateToken middleware - ensures user has a token
func GetOrCreateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip token generation for public endpoints
		publicPaths := []string{"/api/homepage", "/api/movies", "/api/ratings", "/api/health"}
		path := c.Request.URL.Path
		
		isPublic := false
		for _, publicPath := range publicPaths {
			if strings.HasPrefix(path, publicPath) {
				isPublic = true
				break
			}
		}
		
		if !isPublic {
			// Only generate token for user routes
			payload, valid := ValidateToken(c.Request)
			if !valid {
				// Generate new token
				token, err := GenerateAndSetToken(c.Writer)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "Failed to generate user token",
					})
					c.Abort()
					return
				}
				// Store in context for handlers to use
				parts := strings.Split(token, ".")
				payloadJSON, _ := base64.URLEncoding.DecodeString(parts[0])
				var newPayload TokenPayload
				json.Unmarshal(payloadJSON, &newPayload)
				c.Set("user_payload", &newPayload)
			} else {
				c.Set("user_payload", payload)
			}
		}
		c.Next()
	}
}