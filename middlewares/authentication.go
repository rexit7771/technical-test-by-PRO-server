package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"techincal-test/helpers"
	"techincal-test/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		tokenStrings := parts[1]
		claims := &helpers.Claims{}
		token, err := jwt.ParseWithClaims(tokenStrings, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token Signature")
				return
			}
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		if !token.Valid {
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		c.Set("userID", claims.ID)
		c.Set("role", claims.Role)
		c.Set("expiresAt", claims.ExpiresAt.Time)
		c.Next()
	}
}
