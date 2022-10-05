package transaction

import "belajarbwa/user"

type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
