package service_test

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/service"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func setupPaymentService(t *testing.T) (*service.PaymentService, *mockPackageRepo, *mockTransactionRepo) {
	t.Helper()
	pkgRepo := newMockPackageRepo()
	txnRepo := &mockTransactionRepo{}

	duration := 30
	starterPkg := &model.Package{
		ID:            uuid.New(),
		Name:          "Starter",
		Code:          "starter",
		Price:         99000,
		DurationDays:  &duration,
		GuestLimit:    &[]int{100}[0],
		TemplateGroup: "standard",
		Features:      datatypes.JSON(`{"guest.max": 100, "event.max": 3}`),
		IsActive:      true,
	}
	pkgRepo.packages["starter"] = starterPkg

	subSvc := service.NewSubscriptionService(
		&mockSubscriptionRepo{},
		pkgRepo,
		txnRepo,
		newMockUserRepo(),
		&mockAuditLogRepo{},
	)

	paySvc := service.NewPaymentService(
		txnRepo,
		pkgRepo,
		newMockUserRepo(),
		&mockAuditLogRepo{},
		&mockWebhookIdempotencyRepo{},
		config.MidtransConfig{ServerKey: "test-server-key"},
		subSvc,
		noopEmailSender{},
	)
	return paySvc, pkgRepo, txnRepo
}

type noopEmailSender struct{}

func (noopEmailSender) SendAsync(to, subject, htmlBody string) {}

func TestPayment_CreateSnapTransaction(t *testing.T) {
	svc, pkgRepo, _ := setupPaymentService(t)
	ctx := context.Background()
	userID := uuid.New()

	resp, err := svc.CreateSnapTransaction(ctx, userID, dto.CreateSnapRequest{
		PackageID: pkgRepo.packages["starter"].ID.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.OrderID)
	require.NotEmpty(t, resp.SnapToken)
	require.Equal(t, int64(99000), resp.GrossAmount)
	require.Contains(t, resp.SnapRedirectURL, "sandbox.midtrans.com")
}

func TestPayment_CreateSnapTransaction_InvalidPackage(t *testing.T) {
	svc, _, _ := setupPaymentService(t)
	ctx := context.Background()
	userID := uuid.New()

	_, err := svc.CreateSnapTransaction(ctx, userID, dto.CreateSnapRequest{
		PackageID: uuid.New().String(),
	})
	require.Error(t, err)
}

func TestPayment_CreateSnapTransaction_InvalidUUID(t *testing.T) {
	svc, _, _ := setupPaymentService(t)
	ctx := context.Background()
	userID := uuid.New()

	_, err := svc.CreateSnapTransaction(ctx, userID, dto.CreateSnapRequest{
		PackageID: "not-a-uuid",
	})
	require.Error(t, err)
}

func TestPayment_Webhook_SignatureVerification(t *testing.T) {
	svc, _, _ := setupPaymentService(t)

	validPayload := dto.MidtransWebhookPayload{
		OrderID:          "INV-20260101-abc-0001",
		TransactionID:    "txn-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-20260101-abc-0001", "200", "99000"),
	}

	err := svc.HandleWebhook(context.Background(), validPayload)
	require.NoError(t, err)

	invalidPayload := dto.MidtransWebhookPayload{
		OrderID:          "INV-20260101-abc-0001",
		TransactionID:    "txn-002",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     "invalid_signature",
	}

	err = svc.HandleWebhook(context.Background(), invalidPayload)
	require.Error(t, err)
}

func TestPayment_Webhook_DuplicateProtection(t *testing.T) {
	svc, pkgRepo, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-DUPLICATE-001",
		GrossAmount: 99000,
		Status:      "pending",
		UserID:      uuid.New(),
		PackageID:   pkgRepo.packages["starter"].ID,
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-DUPLICATE-001",
		TransactionID:    "txn-dup-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-DUPLICATE-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	err = svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)
}

func TestPayment_Webhook_UnknownOrder(t *testing.T) {
	svc, _, _ := setupPaymentService(t)
	ctx := context.Background()

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-UNKNOWN-001",
		TransactionID:    "txn-unknown-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-UNKNOWN-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)
}

func TestPayment_Webhook_PendingDoesNotActivate(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-PENDING-001",
		GrossAmount: 99000,
		Status:      "pending",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-PENDING-001",
		TransactionID:    "txn-pending-001",
		TransactionStatus: "pending",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-PENDING-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-PENDING-001")
	require.Equal(t, "pending", updatedTxn.Status)
}

func TestPayment_Webhook_ExpireHandled(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-EXPIRE-001",
		GrossAmount: 99000,
		Status:      "pending",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-EXPIRE-001",
		TransactionID:    "txn-expire-001",
		TransactionStatus: "expire",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-EXPIRE-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-EXPIRE-001")
	require.Equal(t, "expire", updatedTxn.Status)
}

func TestPayment_Webhook_CancelHandled(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-CANCEL-001",
		GrossAmount: 99000,
		Status:      "pending",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-CANCEL-001",
		TransactionID:    "txn-cancel-001",
		TransactionStatus: "cancel",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-CANCEL-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-CANCEL-001")
	require.Equal(t, "deny", updatedTxn.Status)
}

func TestPayment_Webhook_DenyHandled(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-DENY-001",
		GrossAmount: 99000,
		Status:      "pending",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-DENY-001",
		TransactionID:    "txn-deny-001",
		TransactionStatus: "deny",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-DENY-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-DENY-001")
	require.Equal(t, "deny", updatedTxn.Status)
}

func TestPayment_Webhook_RefundHandled(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-REFUND-001",
		GrossAmount: 99000,
		Status:      "settlement",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-REFUND-001",
		TransactionID:    "txn-refund-001",
		TransactionStatus: "refund",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-REFUND-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-REFUND-001")
	require.Equal(t, "refund", updatedTxn.Status)
}

func TestPayment_TransactionStateConsistency(t *testing.T) {
	svc, pkgRepo, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-STATE-001",
		GrossAmount: 99000,
		Status:      "pending",
		UserID:      uuid.New(),
		PackageID:   pkgRepo.packages["starter"].ID,
	}
	txnRepo.Create(ctx, txn)

	settlementPayload := dto.MidtransWebhookPayload{
		OrderID:          "INV-STATE-001",
		TransactionID:    "txn-state-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-STATE-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, settlementPayload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-STATE-001")
	require.Equal(t, "settlement", updatedTxn.Status)

	refundPayload := dto.MidtransWebhookPayload{
		OrderID:          "INV-STATE-001",
		TransactionID:    "txn-state-002",
		TransactionStatus: "refund",
		StatusCode:       "201",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-STATE-001", "201", "99000"),
	}

	err = svc.HandleWebhook(ctx, refundPayload)
	require.NoError(t, err)

	updatedTxn, _ = txnRepo.GetByOrderID(ctx, "INV-STATE-001")
	require.Equal(t, "refund", updatedTxn.Status)
}

func TestPayment_Webhook_TimeFields(t *testing.T) {
	svc, pkgRepo, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-TIME-001",
		GrossAmount: 99000,
		Status:      "pending",
		UserID:      uuid.New(),
		PackageID:   pkgRepo.packages["starter"].ID,
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-TIME-001",
		TransactionID:    "txn-time-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		TransactionTime:  "2026-01-15 10:30:00",
		SettlementTime:   "2026-01-15 10:31:00",
		SignatureKey:     generateSignature("INV-TIME-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-TIME-001")
	require.NotNil(t, updatedTxn.TransactionTime)
	require.NotNil(t, updatedTxn.SettlementTime)
	require.Equal(t, "credit_card", updatedTxn.PaymentType)
}

func generateSignature(orderID, statusCode, grossAmount string) string {
	serverKey := "test-server-key"
	hashInput := orderID + statusCode + grossAmount + serverKey
	h := sha512.Sum512([]byte(hashInput))
	return hex.EncodeToString(h[:])
}

func TestPayment_SignatureVerification_Direct(t *testing.T) {
	svc, _, _ := setupPaymentService(t)

	orderID := "INV-SIG-001"
	statusCode := "200"
	grossAmount := "99000"
	signature := generateSignature(orderID, statusCode, grossAmount)

	payload := dto.MidtransWebhookPayload{
		OrderID:          orderID,
		TransactionID:    "txn-sig-001",
		TransactionStatus: "settlement",
		StatusCode:       statusCode,
		GrossAmount:      grossAmount,
		PaymentType:      "credit_card",
		SignatureKey:     signature,
	}

	err := svc.HandleWebhook(context.Background(), payload)
	require.NoError(t, err, "valid signature should be accepted")

	payload.SignatureKey = "wrong_signature"
	err = svc.HandleWebhook(context.Background(), payload)
	require.Error(t, err, "invalid signature should be rejected")
}

func TestPayment_EmptyOrderID(t *testing.T) {
	svc, _, _ := setupPaymentService(t)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "",
		TransactionID:    "txn-empty-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     "any_signature",
	}

	err := svc.HandleWebhook(context.Background(), payload)
	require.Error(t, err, "invalid signature should be rejected")
}

func TestPayment_ListUserTransactions(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		txn := &model.Transaction{
			OrderID:     "INV-LIST-" + string(rune('0'+i)),
			GrossAmount: 99000,
			Status:      "pending",
		}
		txnRepo.Create(ctx, txn)
	}

	txns, total, err := svc.ListUserTransactions(ctx, uuid.Nil, 1, 10)
	require.NoError(t, err)
	require.Len(t, txns, 3)
	require.Equal(t, int64(3), total)
}

func TestPayment_Webhook_AlreadySettledIdempotent(t *testing.T) {
	svc, _, txnRepo := setupPaymentService(t)
	ctx := context.Background()

	txn := &model.Transaction{
		OrderID:     "INV-IDEMPOTENT-001",
		GrossAmount: 99000,
		Status:      "settlement",
	}
	txnRepo.Create(ctx, txn)

	payload := dto.MidtransWebhookPayload{
		OrderID:          "INV-IDEMPOTENT-001",
		TransactionID:    "txn-idempotent-001",
		TransactionStatus: "settlement",
		StatusCode:       "200",
		GrossAmount:      "99000",
		PaymentType:      "credit_card",
		SignatureKey:     generateSignature("INV-IDEMPOTENT-001", "200", "99000"),
	}

	err := svc.HandleWebhook(ctx, payload)
	require.NoError(t, err)

	updatedTxn, _ := txnRepo.GetByOrderID(ctx, "INV-IDEMPOTENT-001")
	require.Equal(t, "settlement", updatedTxn.Status)
}

