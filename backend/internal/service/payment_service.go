package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
)

type EmailSender interface {
	SendAsync(to, subject, htmlBody string)
}

type PaymentService struct {
	txnRepo         repository.TransactionRepository
	pkgRepo         repository.PackageRepository
	userRepo        repository.UserRepository
	auditRepo       repository.AuditLogRepository
	idempotencyRepo repository.WebhookIdempotencyRepository
	midtransKey     string
	subService      *SubscriptionService
	emailSvc        EmailSender
}

func NewPaymentService(txnRepo repository.TransactionRepository, pkgRepo repository.PackageRepository,
	userRepo repository.UserRepository, auditRepo repository.AuditLogRepository,
	idempotencyRepo repository.WebhookIdempotencyRepository,
	cfg config.MidtransConfig, subService *SubscriptionService, emailSvc EmailSender) *PaymentService {
	return &PaymentService{
		txnRepo:           txnRepo,
		pkgRepo:           pkgRepo,
		userRepo:          userRepo,
		auditRepo:         auditRepo,
		idempotencyRepo:   idempotencyRepo,
		midtransKey:       cfg.ServerKey,
		subService:        subService,
		emailSvc:          emailSvc,
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

	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	now := time.Now()
	orderID := fmt.Sprintf("INV-%s-%s-%04d",
		now.Format("20060102"),
		hashUserID(userID.String()),
		now.Unix()%10000)

	transaction := &model.Transaction{
		UserID:      userID,
		PackageID:   pkg.ID,
		OrderID:     orderID,
		GrossAmount: pkg.Price,
		Status:      "pending",
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

	snapToken := orderID + "-" + uuid.New().String()[:8]

	return &dto.SnapResponse{
		TransactionID:   transaction.ID,
		OrderID:         orderID,
		SnapToken:       snapToken,
		SnapRedirectURL: fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v3/redirection/%s", snapToken),
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
	if !s.verifySignature(payload) {
		return errors.ErrSignatureInvalid
	}

	if payload.OrderID == "" {
		return nil
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
		return nil
	}

	internalStatus := mapMidtransStatus(payload.TransactionStatus)

	if txn.Status == "settlement" && internalStatus == "settlement" {
		_ = s.idempotencyRepo.MarkProcessed(ctx, payload.TransactionID, payload.OrderID, internalStatus)
		return nil
	}

	txn.PaymentType = payload.PaymentType
	txn.Status = internalStatus
	if payload.TransactionTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", payload.TransactionTime)
		if err == nil {
			txn.TransactionTime = &t
		}
	}
	if payload.SettlementTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", payload.SettlementTime)
		if err == nil {
			txn.SettlementTime = &t
		}
	}

	if err := s.txnRepo.Update(ctx, txn); err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}

	if internalStatus == "settlement" {
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

func (s *PaymentService) verifySignature(payload dto.MidtransWebhookPayload) bool {
	hashInput := payload.OrderID + payload.StatusCode + payload.GrossAmount + s.midtransKey
	expectedHash := sha512.Sum512([]byte(hashInput))
	expectedSignature := hex.EncodeToString(expectedHash[:])
	return hmac.Equal([]byte(expectedSignature), []byte(payload.SignatureKey))
}

func (s *PaymentService) ListAllTransactions(ctx context.Context, page, perPage int, status string) ([]dto.TransactionResponse, int64, error) {
	txns, total, err := s.txnRepo.ListAll(ctx, page, perPage, status)
	if err != nil {
		return nil, 0, err
	}
	result := make([]dto.TransactionResponse, len(txns))
	for i, txn := range txns {
		result[i] = toTransactionResponse(&txn)
	}
	return result, total, nil
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
