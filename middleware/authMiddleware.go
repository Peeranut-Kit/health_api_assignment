package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware to check if the user is authenticated using JWT
func AuthRequiredMiddleware(c *gin.Context) {
	// Retrieve JWT token from the cookie
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Parse the JWT token
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || claim["staff_hospital_id"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hospital ID not found in token"})
		c.Abort()
		return
	}

	// Convert hospital_id to string
	var hospitalIDStr string
	if id, ok := claim["staff_hospital_id"].(string); ok {
		hospitalIDStr = id
	} else if idFloat, ok := claim["staff_hospital_id"].(float64); ok {
		hospitalIDStr = strconv.FormatFloat(idFloat, 'f', 0, 64) // Convert float64 to string
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid hospital_id format"})
		c.Abort()
		return
	}

	// Set hospital_id in gin.Context
	c.Set("hospital_id", hospitalIDStr)

	// Proceed to the next handler
	c.Next()
}
