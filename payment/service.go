package payment

import (
	"belajarbwa/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-AnPbCSLgJPFR14oh8doclEgV"
	midclient.ClientKey = "SB-Mid-client-kOs1NT3LrL2TPjRy"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
