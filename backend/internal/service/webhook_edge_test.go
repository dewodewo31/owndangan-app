package service_test

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestWebhook_InvalidSignature(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID: "INV-SIG-001", GrossAmount: 99000, Status: "pending",
		UserID: uuid.New(), PackageID: uuid.New(),
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-SIG-001", TransactionID: "txn-sig-001",
		TransactionStatus: "settlement", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: "invalid_signature",
	}

	err := svc.HandleWebhook(ctx, payload)
	require.Error(t, err)
}

func TestWebhook_DuplicateProtection(t *testing.T) {
	svc, pkgRepo, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	userID := uuid.New()
	pkgID := pkgRepo.packages["starter"].ID
	txn := &model.Transaction{
		OrderID: "INV-DUP-001", GrossAmount: 99000, Status: "pending",
		UserID: userID, PackageID: pkgID,
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-DUP-001", TransactionID: "txn-dup-001",
		TransactionStatus: "settlement", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-DUP-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	err = svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-DUP-001")
	require.Equal(t, "settlement", updatedTxn.Status)
}

func TestWebhook_UnknownOrder(t *testing.T) {
	svc, _, _ := setupPaymentService(t)
	ctx := context.Background()

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-UNKNOWN", TransactionID: "txn-unknown",
		TransactionStatus: "settlement", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-UNKNOWN", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)
}

func TestWebhook_PendingDoesNotActivate(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID: "INV-PENDING-001", GrossAmount: 99000, Status: "pending",
		UserID: uuid.New(), PackageID: uuid.New(),
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-PENDING-001", TransactionID: "txn-pending",
		TransactionStatus: "pending", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-PENDING-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-PENDING-001")
	require.Equal(t, "pending", updatedTxn.Status)
}

func TestWebhook_ExpiryHandling(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID: "INV-EXP-001", GrossAmount: 99000, Status: "pending",
		UserID: uuid.New(), PackageID: uuid.New(),
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-EXP-001", TransactionID: "txn-exp",
		TransactionStatus: "expire", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-EXP-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-EXP-001")
	require.Equal(t, "expire", updatedTxn.Status)
}

func TestWebhook_CancelHandling(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID: "INV-CANCEL-001", GrossAmount: 99000, Status: "pending",
		UserID: uuid.New(), PackageID: uuid.New(),
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-CANCEL-001", TransactionID: "txn-cancel",
		TransactionStatus: "cancel", StatusCode: "200",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-CANCEL-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-CANCEL-001")
	require.Equal(t, "deny", updatedTxn.Status)
}

func TestWebhook_RefundHandling(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID: "INV-REF-001", GrossAmount: 99000, Status: "settlement",
		UserID: uuid.New(), PackageID: uuid.New(),
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID: "INV-REF-001", TransactionID: "txn-ref",
		TransactionStatus: "refund", StatusCode: "201",
		GrossAmount: "99000", PaymentType: "credit_card",
		SignatureKey: generateTestSignature("INV-REF-001", "201", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-REF-001")
	require.Equal(t, "refund", updatedTxn.Status)
}

func generateTestSignature(orderID, statusCode, grossAmount string) string {
	serverKey := "test-server-key"
	hashInput := orderID + statusCode + grossAmount + serverKey
	h := sha512.Sum512([]byte(hashInput))
	return hex.EncodeToString(h[:])
}

var _ = datatypes.JSON{}
