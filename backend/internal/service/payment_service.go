package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	idvalidator "github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
)

type EmailSender interface {
	SendAsync(to, subject, htmlBody string)
}

// MidtransClient is the minimal surface of the Midtrans SDK the payment
// service depends on. It keeps the service testable without a second provider
// abstraction and lets tests inject a stub.
type MidtransClient interface {
	CreateSnapTransaction(orderID string, grossAmount int64, customer *midtrans.CustomerDetails, items *[]midtrans.ItemDetails) (*snap.Response, error)
	VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool
}

type PaymentService struct {
	txnRepo         repository.TransactionRepository
	pkgRepo         repository.PackageRepository
	userRepo        repository.UserRepository
	auditRepo       repository.AuditLogRepository
	idempotencyRepo repository.WebhookIdempotencyRepository
	mtClient        MidtransClient
	subService      *SubscriptionService
	emailSvc        EmailSender
}

func NewPaymentService(txnRepo repository.TransactionRepository, pkgRepo repository.PackageRepository,
	userRepo repository.UserRepository, auditRepo repository.AuditLogRepository,
	idempotencyRepo repository.WebhookIdempotencyRepository,
	mtClient MidtransClient, subService *SubscriptionService, emailSvc EmailSender) *PaymentService {
	return &PaymentService{
		txnRepo:         txnRepo,
		pkgRepo:         pkgRepo,
		userRepo:        userRepo,
		auditRepo:       auditRepo,
		idempotencyRepo: idempotencyRepo,
		mtClient:        mtClient,
		subService:      subService,
		emailSvc:        emailSvc,
	}
}

func (s *PaymentService) CreateSnapTransaction(ctx context.Context, userID uuid.UUID, req dto.CreateSnapRequest) (*dto.SnapResponse, error) {
	pkgID, err := uuid.Parse(req.PackageID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid package_id", errors.ErrInvalidInput)
	}

	pkg, err := s.pkgRepo.GetByID(ctx, pkgID)
	if err != nil || pkg == nil {
		return nil, fmt.Errorf("%w: package not found", errors.ErrNotFound)
	}
	if !pkg.IsActive {
		return nil, fmt.Errorf("%w: package is not available for purchase", errors.ErrConflict)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound
	}
	if !idvalidator.IsValidIDPhone(user.Phone) {
		return nil, errors.ErrPhoneRequired
	}

	// Every checkout mints a fresh Snap Token from Midtrans. A previously stored
	// token is never reused: it may have expired or been invalidated and Midtrans
	// would reject it at snap.pay time. A unique order_id per request keeps the
	// DB unique constraint satisfied across repeat checkouts.
	//
	// Midtrans is called BEFORE any row is persisted, so a failed call leaves no
	// dangling pending transaction with an empty snap_token behind.
	now := time.Now()
	orderID := fmt.Sprintf("INV-%s-%s-%d",
		now.Format("20060102"),
		hashUserID(userID.String()),
		now.UnixNano())

	customer := &midtrans.CustomerDetails{
		FName: user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
	items := &[]midtrans.ItemDetails{
		{
			ID:    pkg.Code,
			Name:  pkg.Name,
			Price: pkg.Price,
			Qty:   1,
		},
	}

	snapResp, err := s.mtClient.CreateSnapTransaction(orderID, pkg.Price, customer, items)
	if err != nil {
		return nil, fmt.Errorf("create midtrans snap transaction: %w", err)
	}
	if snapResp == nil || snapResp.Token == "" {
		return nil, fmt.Errorf("create midtrans snap transaction: empty snap token returned")
	}

	transaction := &model.Transaction{
		UserID:      userID,
		PackageID:   pkg.ID,
		OrderID:     orderID,
		GrossAmount: pkg.Price,
		Status:      "pending",
		SnapToken:   snapResp.Token,
	}
	if err := s.txnRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "transaction.created",
		EntityType: "transaction",
		EntityID:   &transaction.ID,
		Metadata:   datatypesJSON(map[string]interface{}{"order_id": orderID, "package_id": pkg.ID.String()}),
	})

	return &dto.SnapResponse{
		TransactionID:   transaction.ID,
		OrderID:         orderID,
		SnapToken:       snapResp.Token,
		SnapRedirectURL: snapResp.RedirectURL,
		GrossAmount:     pkg.Price,
	}, nil
}

func (s *PaymentService) ListUserTransactions(ctx context.Context, userID uuid.UUID, page, perPage int) ([]dto.TransactionResponse, int64, error) {
	txns, total, err := s.txnRepo.ListByUserID(ctx, userID, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	result := make([]dto.TransactionResponse, len(txns))
	for i, txn := range txns {
		result[i] = toTransactionResponse(&txn)
	}
	return result, total, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, payload dto.MidtransWebhookPayload) error {
	if !s.mtClient.VerifySignature(payload.OrderID, payload.StatusCode, payload.GrossAmount, payload.SignatureKey) {
		return errors.ErrSignatureInvalid
	}

	if payload.OrderID == "" {
		return errors.ErrSignatureInvalid
	}

	processed, err := s.idempotencyRepo.IsProcessed(ctx, payload.TransactionID)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	txn, err := s.txnRepo.GetByOrderID(ctx, payload.OrderID)
	if err != nil || txn == nil {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			Action:     "webhook.unknown_order",
			EntityType: "transaction",
			Metadata:   datatypesJSON(map[string]interface{}{"order_id": payload.OrderID}),
		})
		_ = s.idempotencyRepo.MarkProcessed(ctx, payload.TransactionID, payload.OrderID, "unknown")
		return nil
	}

	// Validate the amount reported by Midtrans matches the stored transaction.
	// The server (DB price) is the source of truth; this rejects tampered or
	// forged webhook bodies carrying a different amount.
	notifiedAmount, err := parseGrossAmount(payload.GrossAmount)
	if err != nil || notifiedAmount != txn.GrossAmount {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			Action:     "webhook.amount_mismatch",
			EntityType: "transaction",
			EntityID:   &txn.ID,
			Metadata:   datatypesJSON(map[string]interface{}{"order_id": payload.OrderID, "expected": txn.GrossAmount}),
		})
		_ = s.idempotencyRepo.MarkProcessed(ctx, payload.TransactionID, payload.OrderID, "amount_mismatch")
		return errors.ErrInvalidInput
	}

	internalStatus := mapMidtransStatus(payload.TransactionStatus)

	// State machine: never allow a downgrade to an earlier (lower-rank) status.
	// This blocks late, out-of-order notifications (e.g. expire/deny arriving
	// after settlement) from corrupting a completed transaction.
	currentRank := statusRank[txn.Status]
	newRank := statusRank[internalStatus]
	if newRank < currentRank {
		_ = s.idempotencyRepo.MarkProcessed(ctx, payload.TransactionID, payload.OrderID, internalStatus)
		return nil
	}

	prevStatus := txn.Status
	txn.PaymentType = payload.PaymentType
	txn.Status = internalStatus
	if payload.TransactionTime != "" {
		t, perr := time.Parse("2006-01-02 15:04:05", payload.TransactionTime)
		if perr == nil {
			txn.TransactionTime = &t
		}
	}
	if payload.SettlementTime != "" {
		t, perr := time.Parse("2006-01-02 15:04:05", payload.SettlementTime)
		if perr == nil {
			txn.SettlementTime = &t
		}
	}

	if err := s.txnRepo.Update(ctx, txn); err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}

	if internalStatus == "settlement" && prevStatus != "settlement" {
		sub, err := s.subService.ActivateOnSettlement(ctx, txn.ID)
		if err != nil {
			return fmt.Errorf("activate subscription: %w", err)
		}
		s.sendPaymentSuccessEmail(ctx, txn, sub)
	}

	_ = s.idempotencyRepo.MarkProcessed(ctx, payload.TransactionID, payload.OrderID, internalStatus)

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &txn.UserID,
		Action:     "webhook.processed",
		EntityType: "transaction",
		EntityID:   &txn.ID,
		Metadata:   datatypesJSON(map[string]interface{}{"order_id": payload.OrderID, "status": internalStatus}),
	})

	return nil
}

func (s *PaymentService) ListAllTransactions(ctx context.Context, page, perPage int, status string) ([]dto.TransactionResponse, int64, error) {
	txns, total, err := s.txnRepo.ListAll(ctx, page, perPage, status, "")
	if err != nil {
		return nil, 0, err
	}
	result := make([]dto.TransactionResponse, len(txns))
	for i, txn := range txns {
		result[i] = toTransactionResponse(&txn)
	}
	return result, total, nil
}

// statusRank encodes the lifecycle order of internal transaction statuses.
// Higher rank = later/higher-priority state. Used to reject downgrade
// notifications (e.g. an expire arriving after settlement).
var statusRank = map[string]int{
	"pending":    1,
	"deny":       2,
	"cancel":     2,
	"expire":     2,
	"failed":     2,
	"settlement": 10,
	"refund":     11,
}

func mapMidtransStatus(status string) string {
	switch strings.ToLower(status) {
	case "settlement", "capture":
		return "settlement"
	case "pending":
		return "pending"
	case "deny", "cancel":
		return "deny"
	case "expire":
		return "expire"
	case "failure":
		return "failed"
	case "refund", "partial_refund":
		return "refund"
	default:
		return status
	}
}

func parseGrossAmount(s string) (int64, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(f)), nil
}

func buildRedirectURL(token string) string {
	if token == "" {
		return ""
	}
	return fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v3/redirection/%s", token)
}

func hashUserID(id string) string {
	h := sha512.Sum512([]byte(id))
	return hex.EncodeToString(h[:])[:4]
}

// sendPaymentSuccessEmail is best-effort: any lookup/render failure is dropped,
// delivery happens async via SendAsync — never fails the webhook.
func (s *PaymentService) sendPaymentSuccessEmail(ctx context.Context, txn *model.Transaction, sub *model.Subscription) {
	if s.emailSvc == nil || sub == nil {
		return
	}
	user, err := s.userRepo.GetByID(ctx, txn.UserID)
	if err != nil || user == nil {
		return
	}
	pkg, err := s.pkgRepo.GetByID(ctx, txn.PackageID)
	if err != nil || pkg == nil {
		return
	}
	html, err := email.RenderPaymentSuccess(email.PaymentSuccessData{
		Name:       user.Name,
		PlanName:   pkg.Name,
		Amount:     formatIDR(txn.GrossAmount),
		ExpiryDate: sub.ExpiresAt.Format("2 January 2006"),
	})
	if err != nil {
		return
	}
	s.emailSvc.SendAsync(user.Email, fmt.Sprintf("Pembayaran berhasil — %s aktif!", pkg.Name), html)
}

func formatIDR(amount int64) string {
	s := strconv.FormatInt(amount, 10)
	var out []string
	for len(s) > 3 {
		out = append([]string{s[len(s)-3:]}, out...)
		s = s[:len(s)-3]
	}
	out = append([]string{s}, out...)
	return strings.Join(out, ".")
}

func toTransactionResponse(txn *model.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:             txn.ID,
		OrderID:        txn.OrderID,
		Package:        toPackageBrief(&txn.Package),
		GrossAmount:    txn.GrossAmount,
		Status:         txn.Status,
		PaymentType:    txn.PaymentType,
		TransactionAt:  txn.TransactionTime,
		SettlementTime: txn.SettlementTime,
		CreatedAt:      txn.CreatedAt,
	}
}
