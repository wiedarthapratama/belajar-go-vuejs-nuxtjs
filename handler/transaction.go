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
