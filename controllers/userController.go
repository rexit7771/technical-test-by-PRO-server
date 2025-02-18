package controllers

import (
	"fmt"
	"net/http"
	"techincal-test/database"
	"techincal-test/helpers"
	"techincal-test/response"
	"techincal-test/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var err error

func AddNewUser(c *gin.Context) {
	var newUser structs.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := newUser.Validate(); err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
	if err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newUser.Password = string(hashedPassword)

	if err := database.DB.Create(&newUser).Error; err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.CommonResponse(c, http.StatusCreated, "User registered successfully")
}

func Login(c *gin.Context) {
	var user structs.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user.Email == "" {
		response.AbortResponse(c, http.StatusBadRequest, "Email is required")
		return
	}

	if user.Password == "" {
		response.AbortResponse(c, http.StatusBadRequest, "Password is required")
		return
	}

	var userDB structs.User
	tx := database.DB.Where("email = ?", user.Email).First(&userDB)
	if tx.Error != nil {
		response.AbortResponse(c, http.StatusNotFound, "Invalid Email / Password")
		return
	}

	bcryptResult := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if bcryptResult != nil {
		response.AbortResponse(c, http.StatusNotFound, "Invalid Email / Password")
		return
	}

	token, tokenError := helpers.SignPayLoad(userDB)
	if tokenError != nil {
		fmt.Print(tokenError)
		response.AbortResponse(c, http.StatusInternalServerError, tokenError.Error())
		return
	}

	response.TokenResponse(c, token)
}

func GetUserByIdAuth(c *gin.Context) {
	userID, exists := c.Get("userID")
	if exists == false {
		response.AbortResponse(c, http.StatusUnauthorized, "You need to log in first")
		return
	}
	userIDUint := userID.(uint)

	var user structs.User
	err := database.DB.First(&user, userIDUint).Error
	if err != nil {

		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.CommonResponse(c, http.StatusOK, user)
}

func GetUserByIdParam(c *gin.Context) {
	userID := c.Param("id")
	var user structs.User
	err := database.DB.First(&user, userID).Error
	if err != nil {
		message := "User with id " + userID + " is not found"
		response.AbortResponse(c, http.StatusNotFound, message)
		return
	}

	response.CommonResponse(c, http.StatusOK, user)
}

func GetAllUser(c *gin.Context) {
	var users []structs.User
	query := database.DB.Model(&structs.User{})
	helpers.QuerySearch(c, query)
	err := query.Find(&users).Error
	if err != nil {
		response.AbortResponse(c, http.StatusNotFound, "Users data is not found")
		return
	}

	response.CommonResponse(c, http.StatusOK, users)
}

func EditUserById(c *gin.Context) {
	userID := c.Param("id")
	var userDB structs.User
	err := database.DB.First(&userDB, userID).Error
	if err != nil {
		message := "User with id " + userID + " is not found"
		response.AbortResponse(c, http.StatusNotFound, message)
		return
	}

	var userUpdate structs.User
	if err := c.ShouldBindBodyWithJSON(&userUpdate); err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userUpdate.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.MinCost)
		if err != nil {
			fmt.Print(err.Error())
			response.AbortResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		userUpdate.Password = string(hashedPassword)
	} else {
		userUpdate.Password = userDB.Password
	}

	err = userUpdate.Validate()
	if err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = database.DB.Model(&userDB).Updates(userUpdate).Error
	if err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message := "User with id " + userID + " has been updated"
	response.CommonResponse(c, http.StatusOK, message)
}

func DeleteUserById(c *gin.Context) {
	userID := c.Param("id")
	var userDB structs.User
	err := database.DB.First(&userDB, userID).Error
	if err != nil {
		message := "User with id " + userID + " is not found"
		response.AbortResponse(c, http.StatusNotFound, message)
		return
	}

	err = database.DB.Unscoped().Delete(&userDB).Error
	if err != nil {
		fmt.Print(err.Error())
		response.AbortResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := "User with id " + userID + " has been deleted"
	response.CommonResponse(c, http.StatusOK, message)
}
