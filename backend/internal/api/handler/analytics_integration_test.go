package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestAnalytics_TrackEvent_Published(t *testing.T) {
	server, db, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Event")
	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/analytics/events", map[string]string{
		"event_id": eventID,
		"type":     "whatsapp_click",
	}, "")
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var row model.AnalyticsEvent
	require.NoError(t, db.Where("event_id = ? AND event_type = ?", eventID, "whatsapp_click").First(&row).Error)
	require.NotEmpty(t, row.IPAddress)
}

func TestAnalytics_TrackEvent_InvalidType(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Invalid")
	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/analytics/events", map[string]string{
		"event_id": eventID,
		"type":     "bogus_click",
	}, "")
	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
}

func TestAnalytics_TrackEvent_UnpublishedNotFound(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Unpublished")

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/analytics/events", map[string]string{
		"event_id": eventID,
		"type":     "map_click",
	}, "")
	require.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}

func TestAnalytics_GetEventAnalytics(t *testing.T) {
	server, db, userID, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Stats")
	eventUUID := uuid.MustParse(eventID)

	analyticsRepo := repository.NewAnalyticsEventRepository(db)
	ips := []string{"1.1.1.1", "1.1.1.1", "2.2.2.2"}
	for _, ip := range ips {
		require.NoError(t, analyticsRepo.Create(t.Context(), &model.AnalyticsEvent{
			EventID:   &eventUUID,
			EventType: "page_view",
			IPAddress: ip,
		}))
	}
	require.NoError(t, analyticsRepo.Create(t.Context(), &model.AnalyticsEvent{
		EventID:   &eventUUID,
		EventType: "whatsapp_click",
		IPAddress: "3.3.3.3",
	}))

	// Analytics is a paid feature: upgrade to starter before reading stats.
	starter := &model.Package{}
	require.NoError(t, db.Where("code = ?", "starter").First(starter).Error)
	require.NoError(t, db.Model(&model.Subscription{}).Where("user_id = ?", userID).Update("package_id", starter.ID).Error)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/analytics", nil, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	require.Equal(t, float64(3), data["views"])
	require.Equal(t, float64(2), data["unique_views"])
	require.Equal(t, float64(1), data["whatsapp_clicks"])
	require.Equal(t, float64(0), data["map_clicks"])
	require.Equal(t, float64(0), data["phone_clicks"])
	require.Equal(t, float64(0), data["rsvp_count"])
}

func TestAnalytics_GetEventAnalytics_Entitlement(t *testing.T) {
	server, db, userID, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Entitlement")

	// Free plan: analytics disabled -> 403
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/analytics", nil, token)
	require.Equal(t, http.StatusForbidden, w.Code, w.Body.String())

	// Upgrade to starter: analytics enabled -> 200
	starter := &model.Package{}
	require.NoError(t, db.Where("code = ?", "starter").First(starter).Error)
	require.NoError(t, db.Model(&model.Subscription{}).Where("user_id = ?", userID).Update("package_id", starter.ID).Error)

	w2 := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/analytics", nil, token)
	require.Equal(t, http.StatusOK, w2.Code, w2.Body.String())
}

func TestAnalytics_GetEventAnalytics_Ownership(t *testing.T) {
	server, db, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Analytics Ownership")

	otherUser := &model.User{Name: "Other", Email: "other-" + uuid.New().String() + "@test.com", PasswordHash: "x", Role: "user", Status: "active"}
	require.NoError(t, db.Create(otherUser).Error)
	freePkg := &model.Package{}
	require.NoError(t, db.Where("code = ?", "free").First(freePkg).Error)
	expiry := time.Now().Add(90 * 24 * time.Hour)
	require.NoError(t, db.Create(&model.Subscription{UserID: otherUser.ID, PackageID: freePkg.ID, Status: "active", StartAt: time.Now(), ExpiresAt: &expiry}).Error)
	otherToken := tokenForUser(t, otherUser.ID)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/analytics", nil, otherToken)
	require.Equal(t, http.StatusForbidden, w.Code, w.Body.String())
}
