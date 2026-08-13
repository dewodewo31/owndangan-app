package midtrans

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Client struct {
	snapClient   *snap.Client
	coreClient   *coreapi.Client
	serverKey    string
	isProduction bool
}

func NewClient(serverKey string, isProduction bool) *Client {
	s := snap.Client{}
	s.New(serverKey, envType(isProduction))

	c := coreapi.Client{}
	c.New(serverKey, envType(isProduction))

	return &Client{
		snapClient:   &s,
		coreClient:   &c,
		serverKey:    serverKey,
		isProduction: isProduction,
	}
}

func envType(isProduction bool) midtrans.EnvironmentType {
	if isProduction {
		return midtrans.Production
	}
	return midtrans.Sandbox
}

func (c *Client) CreateSnapTransaction(orderID string, grossAmount int64, customer *midtrans.CustomerDetails, items *[]midtrans.ItemDetails) (*snap.Response, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		CustomerDetail: customer,
		Items:          items,
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeCreditCard,
			snap.PaymentTypeBCAVA,
			snap.PaymentTypeBNIVA,
			snap.PaymentTypeBRIVA,
			snap.PaymentTypePermataVA,
			snap.PaymentTypeOtherVA,
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeIndomaret,
			snap.PaymentTypeAlfamart,
			snap.PaymentTypeAkulaku,
		},
	}

	resp, err := c.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("create snap transaction: %w", err)
	}
	return resp, nil
}

func (c *Client) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	hashInput := orderID + statusCode + grossAmount + c.serverKey
	expected := sha512Sum(hashInput)
	return hmacEqual(expected, signatureKey)
}

func (c *Client) CheckTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	resp, err := c.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("check transaction: %w", err)
	}
	return resp, nil
}

func (c *Client) IsProduction() bool {
	return c.isProduction
}

func sha512Sum(input string) string {
	h := sha512.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func hmacEqual(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}
