package transaction

import (
	"belajarbwa/campaign"
	"belajarbwa/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Campaign   campaign.Campaign
}
