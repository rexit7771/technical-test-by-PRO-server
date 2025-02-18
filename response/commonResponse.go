package response

import "github.com/gin-gonic/gin"

func CommonResponse(c *gin.Context, status int, result interface{}) {
	c.JSON(status, gin.H{"result": result})
}
