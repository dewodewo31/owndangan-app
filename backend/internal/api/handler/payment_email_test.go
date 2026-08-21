package handler_test

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type recEmail struct{ To, Subject, HTML string }

type recEmailSender struct {
	mu     sync.Mutex
	emails []recEmail
}

func (r *recEmailSender) SendAsync(to, subject, htmlBody string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.emails = append(r.emails, recEmail{To: to, Subject: subject, HTML: htmlBody})
}

func (r *recEmailSender) SendWithRetry(to, subject, htmlBody string) error {
	r.SendAsync(to, subject, htmlBody)
	return nil
}

func (r *recEmailSender) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.emails)
}

func midtransSignature(orderID, statusCode, grossAmount, key string) string {
	h := sha512.Sum512([]byte(orderID + statusCode + grossAmount + key))
	return hex.EncodeToString(h[:])
}

func registerPayer(t *testing.T) (string, string) {
	testCounter++
	email := fmt.Sprintf("payer_%d@example.com", testCounter)
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Payer", "email": email, "password": "securepassword123",
	}, "")
	lr := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusOK, lr.Code, lr.Body.String())
	var ld map[string]interface{}
	require.NoError(t, json.Unmarshal(lr.Body.Bytes(), &ld))
	token := ld["data"].(map[string]interface{})["access_token"].(string)
	return token, email
}

func starterPackageID(t *testing.T) string {
	w := doAuthRequest(t, http.MethodGet, "/api/v1/packages", nil, "")
	require.Equal(t, http.StatusOK, w.Code)
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	for _, p := range data["data"].([]interface{}) {
		pkg := p.(map[string]interface{})
		if pkg["code"].(string) == "starter" {
			return pkg["id"].(string)
		}
	}
	t.Fatal("starter package not found")
	return ""
}

func createSnapOrder(t *testing.T, token, pkgID string) string {
	w := doAuthRequest(t, http.MethodPost, "/api/v1/payments/snap", map[string]string{
		"package_id": pkgID,
	}, token)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	return data["data"].(map[string]interface{})["order_id"].(string)
}

func sendWebhook(t *testing.T, orderID, txnID, status, gross string) {
	w := doAuthRequest(t, http.MethodPost, "/api/v1/payments/webhook", map[string]string{
		"order_id":           orderID,
		"transaction_id":     txnID,
		"status_code":        "200",
		"gross_amount":       gross,
		"transaction_status": status,
		"payment_type":       "bank_transfer",
		"signature_key":      midtransSignature(orderID, "200", gross, ""),
	}, "")
	require.Equal(t, http.StatusNoContent, w.Code, w.Body.String())
}

func TestPaymentWebhook_Settlement_SendsEmailOnce(t *testing.T) {
	rec := &recEmailSender{}
	testEmailSender = rec
	defer func() { testEmailSender = nil }()
	setupAuthTestServer(t)

	token, payerEmail := registerPayer(t)
	orderID := createSnapOrder(t, token, starterPackageID(t))

	sendWebhook(t, orderID, "TXN-SETTLE-1", "settlement", "99000.00")
	require.Equal(t, 1, rec.count(), "settlement must trigger exactly one email")

	sendWebhook(t, orderID, "TXN-SETTLE-1", "settlement", "99000.00")
	require.Equal(t, 1, rec.count(), "duplicate webhook must not double-send")

	rec.mu.Lock()
	defer rec.mu.Unlock()
	require.Equal(t, payerEmail, rec.emails[0].To)
	require.Contains(t, rec.emails[0].Subject, "aktif")
	require.Contains(t, rec.emails[0].HTML, "99.000")
}

func TestPaymentWebhook_Deny_NoEmail(t *testing.T) {
	rec := &recEmailSender{}
	testEmailSender = rec
	defer func() { testEmailSender = nil }()
	setupAuthTestServer(t)

	token, _ := registerPayer(t)
	orderID := createSnapOrder(t, token, starterPackageID(t))

	sendWebhook(t, orderID, "TXN-DENY-1", "deny", "99000.00")
	require.Equal(t, 0, rec.count(), "deny must not send any email")
}
