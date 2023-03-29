package handler

import (
	"belajarbwa/helper"
	"belajarbwa/transaction"
	"belajarbwa/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to get campaign transactions", http.StatusBadRequest, "error", err)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// mendapatkan user sesuai token
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to get campaign transactionss", http.StatusBadRequest, "error", err)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)

	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to get users transactionss", http.StatusBadRequest, "error", err)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mendapatkan user sesuai token
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create Campaign", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to process notification", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.ProcessPayment(input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to process notification 2", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
