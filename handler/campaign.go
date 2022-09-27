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
