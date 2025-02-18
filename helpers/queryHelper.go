package helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func QuerySearch(c *gin.Context, query *gorm.DB) {
	searchName := c.DefaultQuery("name", "")
	searchEmail := c.DefaultQuery("email", "")

	if searchName != "" {
		queryName := "%" + searchName + "%"
		query = query.Where("name ILIKE ?", queryName)
	}

	if searchEmail != "" {
		queryEmail := "%" + searchEmail + "%"
		query = query.Where("email ILIKE ?", queryEmail)
	}
}
