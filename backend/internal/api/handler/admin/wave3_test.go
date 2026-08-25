package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/cache"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"

	"github.com/google/uuid"
)

func wave3RegisterLogin(t *testing.T, email string) (string, uuid.UUID) {
	t.Helper()
	doAdminRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": email, "email": email, "password": "password123",
	}, "")
	loginResp := doAdminRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "password123",
	}, "")
	var loginData map[string]interface{}
	require.NoError(t, json.Unmarshal(loginResp.Body.Bytes(), &loginData))
	data := loginData["data"].(map[string]interface{})
	token := data["access_token"].(string)
	userID, _ := uuid.Parse(data["id"].(string))
	return token, userID
}

func wave3CreateEvent(t *testing.T, token string) (string, string) {
	t.Helper()
	w := doAdminRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "Wave3 Wedding", "groom_name": "A", "bride_name": "B", "wedding_date": "2026-09-01",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	return data["id"].(string), data["slug"].(string)
}

func wave3CreateGuest(t *testing.T, token, eventID string) (string, string) {
	t.Helper()
	w := doAdminRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "Guest One", "phone": "62812",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	return data["id"].(string), data["token"].(string)
}

func wave3UpgradeToPro(t *testing.T, userID uuid.UUID) {
	t.Helper()
	db := openAdminTestDB(t)
	db.Model(&model.Subscription{}).Where("user_id = ?", userID).Update("status", "canceled")
	var pkg model.Package
	db.Where("code = ?", "premium").First(&pkg)
	start := time.Now()
	expires := start.Add(30 * 24 * time.Hour)
	sub := &model.Subscription{UserID: userID, PackageID: pkg.ID, Status: "active", StartAt: start, ExpiresAt: &expires}
	require.NoError(t, db.Create(sub).Error)
}

func TestWave3_RSVPExport_CSV(t *testing.T) {
	setupAdminTestServer(t)
	token, _ := wave3RegisterLogin(t, "export_csv@example.com")
	eventID, _ := wave3CreateEvent(t, token)
	doAdminRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	_, gToken := wave3CreateGuest(t, token, eventID)
	doAdminRequest(t, http.MethodPost, "/api/v1/rsvp/"+eventID+"/submit", map[string]interface{}{
		"token": gToken, "attendance": "attending", "guest_count": 2,
	}, "")

	w := doAdminRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/rsvp/export?format=csv", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	require.Contains(t, w.Header().Get("Content-Disposition"), "rsvp-")
	require.Contains(t, w.Body.String(), "Guest One")
}

func TestWave3_RSVPExport_XLSX_Pro(t *testing.T) {
	setupAdminTestServer(t)
	token, userID := wave3RegisterLogin(t, "export_xlsx@example.com")
	wave3UpgradeToPro(t, userID)
	eventID, _ := wave3CreateEvent(t, token)
	_, gToken := wave3CreateGuest(t, token, eventID)
	doAdminRequest(t, http.MethodPost, "/api/v1/rsvp/"+eventID+"/submit", map[string]interface{}{
		"token": gToken, "attendance": "attending", "guest_count": 2,
	}, "")

	w := doAdminRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/rsvp/export?format=xlsx", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Header().Get("Content-Type"), "spreadsheetml")
	require.Contains(t, w.Header().Get("Content-Disposition"), ".xlsx")

	f, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	require.NoError(t, err)
	require.NotNil(t, f)
}

func TestWave3_RSVPExport_XLSX_NonPro(t *testing.T) {
	setupAdminTestServer(t)
	token, _ := wave3RegisterLogin(t, "export_nonpro@example.com")
	eventID, _ := wave3CreateEvent(t, token)

	w := doAdminRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/rsvp/export?format=xlsx", nil, token)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestWave3_QRCheckIn_Success(t *testing.T) {
	setupAdminTestServer(t)
	token, userID := wave3RegisterLogin(t, "checkin_ok@example.com")
	wave3UpgradeToPro(t, userID)
	eventID, _ := wave3CreateEvent(t, token)
	_, gToken := wave3CreateGuest(t, token, eventID)

	w := doAdminRequest(t, http.MethodPost, "/api/v1/guests/check-in", map[string]string{"token": gToken}, "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	require.NotEmpty(t, data["attended_at"])
}

func TestWave3_QRCheckIn_Idempotent(t *testing.T) {
	setupAdminTestServer(t)
	token, userID := wave3RegisterLogin(t, "checkin_idem@example.com")
	wave3UpgradeToPro(t, userID)
	eventID, _ := wave3CreateEvent(t, token)
	_, gToken := wave3CreateGuest(t, token, eventID)

	w1 := doAdminRequest(t, http.MethodPost, "/api/v1/guests/check-in", map[string]string{"token": gToken}, "")
	require.Equal(t, http.StatusOK, w1.Code)
	w2 := doAdminRequest(t, http.MethodPost, "/api/v1/guests/check-in", map[string]string{"token": gToken}, "")
	require.Equal(t, http.StatusOK, w2.Code)
}

func TestWave3_QRCheckIn_NonPro(t *testing.T) {
	setupAdminTestServer(t)
	token, _ := wave3RegisterLogin(t, "checkin_nonpro@example.com")
	eventID, _ := wave3CreateEvent(t, token)
	_, gToken := wave3CreateGuest(t, token, eventID)

	w := doAdminRequest(t, http.MethodPost, "/api/v1/guests/check-in", map[string]string{"token": gToken}, "")
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestWave3_Cache_InvalidationOnUpdate(t *testing.T) {
	setupAdminTestServer(t)
	token, _ := wave3RegisterLogin(t, "cache_inv@example.com")
	eventID, slug := wave3CreateEvent(t, token)

	pub := doAdminRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	require.Equal(t, http.StatusOK, pub.Code)

	// cache miss -> load from repo and store
	w1 := doAdminRequest(t, http.MethodGet, "/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, w1.Code)
	require.Contains(t, w1.Body.String(), "Wave3 Wedding")

	// update the invitation -> cache must be invalidated
	upd := doAdminRequest(t, http.MethodPut, "/api/v1/events/"+eventID, map[string]string{"title": "Updated Title"}, token)
	require.Equal(t, http.StatusOK, upd.Code)

	// subsequent read reflects the update (proves invalidation)
	w2 := doAdminRequest(t, http.MethodGet, "/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, w2.Code)
	require.Contains(t, w2.Body.String(), "Updated Title")
}

func TestWave3_PublicCache_Unit(t *testing.T) {
	c := cache.NewTTL(50 * time.Millisecond)

	_, ok := c.Get("missing")
	require.False(t, ok)

	c.Set("k", "v")
	v, ok := c.Get("k")
	require.True(t, ok)
	require.Equal(t, "v", v)

	c.Delete("k")
	_, ok = c.Get("k")
	require.False(t, ok)

	// expired entries are evicted on access
	c.Set("e", "x")
	time.Sleep(60 * time.Millisecond)
	_, ok = c.Get("e")
	require.False(t, ok)
}
