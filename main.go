package main

import (
	"os"
	"techincal-test/controllers"
	"techincal-test/database"
	"techincal-test/middlewares"
	"techincal-test/structs"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	database.DB.AutoMigrate(&structs.User{})

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.POST("/register", controllers.AddNewUser)
	router.POST("/login", controllers.Login)
	router.POST("/login/admin", controllers.LoginAdmin)
	userGroup := router.Group("/users")
	{
		userGroup.Use(middlewares.Authentication())
		userGroup.GET("", controllers.GetUserByIdAuth)
		userGroup.Use(middlewares.IsAdmin())
		userGroup.GET("/all", controllers.GetAllUser)
		userGroup.GET("/:id", controllers.GetUserByIdParam)
		userGroup.PUT("/:id", controllers.EditUserById)
		userGroup.DELETE("/:id", controllers.DeleteUserById)
	}
	router.Run(":" + os.Getenv("PORT"))
}
