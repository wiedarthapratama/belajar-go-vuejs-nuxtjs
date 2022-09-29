package handler

import (
	"belajarbwa/campaign"
	"belajarbwa/helper"
	"belajarbwa/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// string to int
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// format responsenya
	response := helper.APIResponse("List of Campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	// response to json
	c.JSON(http.StatusOK, response)
}

// api/v1/campaigns/1
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Error to get detail of campaign", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Error to get detail of campaign", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to get Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mendapatkan user sesuai token
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create Campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to update Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		// error to format error
		errors := helper.FormatValidationError(err)
		// mapping apa aja ke object errors
		errorMessage := gin.H{"errors": errors}
		// format responsenya
		response := helper.APIResponse("Failed to update Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		// response to json
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mendapatkan user sesuai token
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		// format responsenya
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		// response to json
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create Campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}
