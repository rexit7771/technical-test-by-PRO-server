package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenResponse(c *gin.Context, token string) {
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
