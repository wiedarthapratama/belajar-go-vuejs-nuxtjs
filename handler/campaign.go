package handler

import (
	"belajarbwa/campaign"
	"belajarbwa/helper"
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
