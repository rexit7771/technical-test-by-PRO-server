package response

import (
	"github.com/gin-gonic/gin"
)

func AbortResponse(c *gin.Context, status int, message interface{}) {
	switch message {
	case "ERROR: duplicate key value violates unique constraint \"uni_users_email\" (SQLSTATE 23505)":
		message = "Email has been used, use another email"
	case "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag":
		message = "Invalid Email Format"
	case "Key: 'User.BirthDate' Error:Field validation for 'BirthDate' failed on the 'birthdate' tag":
		message = "Birth Date must be before today"
	}

	c.AbortWithStatusJSON(status, gin.H{"result": message})
}
