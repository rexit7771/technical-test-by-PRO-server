package middlewares

import (
	"net/http"
	"techincal-test/response"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		role, exists := c.Get("role")
		if !exists {
			response.AbortResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}

		if role != "admin" {
			response.AbortResponse(c, http.StatusForbidden, "You're Unauthorized to access this")
			return
		} else if role == "admin" {
			c.Next()
		}
	}
}
