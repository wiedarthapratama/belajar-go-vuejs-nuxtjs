package handler

import (
	"belajarbwa/auth"
	"belajarbwa/helper"
	"belajarbwa/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// inputan dari user
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// fungsi register user
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Make User Failed", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// generate jwt token
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Make Token Failed", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// format user untuk response
	formatter := user.FormatUser(newUser, token)
	// format responsenya
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	// response to json
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	// inputan dari user
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// proses login
	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": err.Error()}
		// format responsenya
		response := helper.APIResponse("Login Failed", http.StatusNotFound, "error", errorMessage)
		// response to json
		c.JSON(http.StatusNotFound, response)
		return
	}
	// generate jwt token
	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Make Token Failed", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// format user untuk response
	formatter := user.FormatUser(loggedinUser, token)
	// format responsenya
	response := helper.APIResponse("Successfuly Loggedin", http.StatusOK, "success", formatter)
	// response to json
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": "Server Error"}
		// format responsenya
		response := helper.APIResponse("Email Checking Failed", http.StatusInternalServerError, "error", errorMessage)
		// response to json
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// bikin objek is available
	data := gin.H{
		"is_available": isEmailAvailable,
	}
	// string message meta
	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is Available"
	}
	// format responsenya
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	// response to json
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		// mapping apa aja ke object errors
		data := gin.H{"is_uploaded": false}
		// format responsenya
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// harusnya dapet dari jwt
	userID := 1
	// path gambar %d = userID %s = file.filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	// upload file
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		// mapping apa aja ke object errors
		data := gin.H{"is_uploaded": false}
		// format responsenya
		response := helper.APIResponse("Failed to store avatar image", http.StatusBadRequest, "error", data)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		// mapping apa aja ke object errors
		data := gin.H{"is_uploaded": false}
		// format responsenya
		response := helper.APIResponse("Failed to update avatar image", http.StatusBadRequest, "error", data)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// mapping apa aja ke object errors
	data := gin.H{"is_uploaded": true}
	// format responsenya
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	// response to json
	c.JSON(http.StatusOK, response)
}
