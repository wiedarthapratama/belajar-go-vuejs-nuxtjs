package handler

import (
	"belajarbwa/helper"
	"belajarbwa/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// inputan dari user
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mappring apa aja ke object errors
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
	// format user untuk response
	formatter := user.FormatUser(newUser, "initokennya")
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
		// mappring apa aja ke object errors
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
		// mappring apa aja ke object errors
		errorMessage := gin.H{"errors": err.Error()}
		// format responsenya
		response := helper.APIResponse("Login Failed", http.StatusNotFound, "error", errorMessage)
		// response to json
		c.JSON(http.StatusNotFound, response)
		return
	}
	// format user untuk response
	formatter := user.FormatUser(loggedinUser, "initokennya")
	// format responsenya
	response := helper.APIResponse("Successfuly Loggedin", http.StatusOK, "success", formatter)
	// response to json
	c.JSON(http.StatusOK, response)
}
