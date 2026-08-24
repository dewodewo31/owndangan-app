package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSnapRequest struct {
	PackageID   string `json:"package_id" validate:"required,uuid"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

type SnapResponse struct {
	TransactionID   uuid.UUID `json:"transaction_id"`
	OrderID         string    `json:"order_id"`
	SnapToken       string    `json:"snap_token"`
	SnapRedirectURL string    `json:"snap_redirect_url,omitempty"`
	GrossAmount     int64     `json:"gross_amount"`
}

type TransactionResponse struct {
	ID             uuid.UUID    `json:"id"`
	OrderID        string       `json:"order_id"`
	Package        PackageBrief `json:"package"`
	GrossAmount    int64        `json:"gross_amount"`
	Status         string       `json:"status"`
	PaymentType    string       `json:"payment_type,omitempty"`
	TransactionAt  *time.Time   `json:"transaction_time,omitempty"`
	SettlementTime *time.Time   `json:"settlement_time,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

type MidtransWebhookPayload struct {
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	SettlementTime    string `json:"settlement_time,omitempty"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	Bank              string `json:"bank,omitempty"`
}
