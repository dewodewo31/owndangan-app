package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func tokenForUser(t *testing.T, userID uuid.UUID) string {
	t.Helper()
	jwtSvc := jwt.New("test-secret", 15*time.Minute, 7*24*time.Hour)
	token, _, err := jwtSvc.GenerateAccessToken(userID, "user")
	require.NoError(t, err)
	return token
}

func setupEditorTest(t *testing.T) (*api.Server, *gorm.DB, uuid.UUID, string) {
	t.Helper()
	server, db := setupTestServer(t)

	user := &model.User{
		Name:         "Editor User",
		Email:        "editor-" + uuid.New().String() + "@test.com",
		PasswordHash: "irrelevant",
		Role:         "user",
		Status:       "active",
	}
	require.NoError(t, db.Create(user).Error)

	freePkg := &model.Package{}
	require.NoError(t, db.Where("code = ?", "free").First(freePkg).Error)

	expiry := time.Now().Add(90 * 24 * time.Hour)
	sub := &model.Subscription{
		UserID:    user.ID,
		PackageID: freePkg.ID,
		Status:    "active",
		StartAt:   time.Now(),
		ExpiresAt: &expiry,
	}
	require.NoError(t, db.Create(sub).Error)

	return server, db, user.ID, tokenForUser(t, user.ID)
}

func createTestEvent(t *testing.T, server *api.Server, token, title string) string {
	t.Helper()
	body := map[string]string{
		"title":        title,
		"couple_name":  title,
		"groom_name":   "Adi",
		"bride_name":   "Aisyah",
		"wedding_date": "2025-12-31",
	}
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events", body, token)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	return data["id"].(string)
}

func slugOf(t *testing.T, server *api.Server, token, eventID string) string {
	t.Helper()
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	var gresp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &gresp))
	return gresp["data"].(map[string]interface{})["slug"].(string)
}

func TestEditor_UpdateBasicFields_Persists(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Editor Event")

	patch := map[string]string{"bride_name": "Aisyah"}
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID, patch, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w2 := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	require.Equal(t, http.StatusOK, w2.Code, w2.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	require.Equal(t, "Aisyah", data["bride_name"].(string))
}

func TestEditor_PublicView_RendersEditedData(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Public Event")

	patch := map[string]string{"bride_name": "Bunga", "groom_name": "Adi"}
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID, patch, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)
	slug := slugOf(t, server, token, eventID)

	wPublic := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, wPublic.Code, wPublic.Body.String())
	var pub map[string]interface{}
	require.NoError(t, json.Unmarshal(wPublic.Body.Bytes(), &pub))
	ev := pub["data"].(map[string]interface{})["event"].(map[string]interface{})
	require.Equal(t, "Bunga", ev["bride_name"].(string))
	require.Equal(t, "Adi", ev["groom_name"].(string))
}

func TestEditor_SectionsToggles_PersistAndPublic(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Sections Event")

	sections := map[string]interface{}{
		"gallery_enabled":       false,
		"guestbook_enabled":     false,
		"digital_gifts_enabled": true,
		"dress_code":            "Black tie",
		"closing_message":       "See you there!",
	}
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/sections", sections, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	wg := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/sections", nil, token)
	require.Equal(t, http.StatusOK, wg.Code, wg.Body.String())
	var sec map[string]interface{}
	require.NoError(t, json.Unmarshal(wg.Body.Bytes(), &sec))
	s := sec["data"].(map[string]interface{})
	require.Equal(t, false, s["gallery_enabled"].(bool))
	require.Equal(t, false, s["guestbook_enabled"].(bool))
	require.Equal(t, true, s["digital_gifts_enabled"].(bool))
	require.Equal(t, "Black tie", s["dress_code"].(string))

	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)
	slug := slugOf(t, server, token, eventID)

	wp := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, wp.Code, wp.Body.String())
	var pub map[string]interface{}
	require.NoError(t, json.Unmarshal(wp.Body.Bytes(), &pub))
	sectionsDTO := pub["data"].(map[string]interface{})["sections"].(map[string]interface{})
	require.Equal(t, false, sectionsDTO["gallery_enabled"].(bool))
	require.Equal(t, false, sectionsDTO["guestbook_enabled"].(bool))
	require.Equal(t, true, sectionsDTO["digital_gifts_enabled"].(bool))
}

func TestEditor_DigitalGifts_PersistsAndPublic(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Gift Event")

	gift := map[string]interface{}{
		"bank_accounts": []map[string]interface{}{{"bank": "Mandiri", "account": "123", "name": "Budi"}},
		"ewallet":       map[string]interface{}{"ovo": "0812"},
		"gift_message":  "Terima kasih",
	}
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/digital-gifts", gift, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	wg := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/digital-gifts", nil, token)
	require.Equal(t, http.StatusOK, wg.Code, wg.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(wg.Body.Bytes(), &resp))
	g := resp["data"].(map[string]interface{})
	require.Equal(t, "Terima kasih", g["gift_message"].(string))
	require.Len(t, g["bank_accounts"].([]interface{}), 1)

	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)
	slug := slugOf(t, server, token, eventID)

	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/sections",
		map[string]interface{}{"digital_gifts_enabled": true}, token).Code)

	wp := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, wp.Code, wp.Body.String())
	var pub map[string]interface{}
	require.NoError(t, json.Unmarshal(wp.Body.Bytes(), &pub))
	dg := pub["data"].(map[string]interface{})["digital_gift"].(map[string]interface{})
	require.Equal(t, "Terima kasih", dg["gift_message"].(string))
}

func uploadGalleryFile(t *testing.T, server *api.Server, token, eventID, filename string) string {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	part, err := mw.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write([]byte("fake-image-bytes"))
	require.NoError(t, err)
	require.NoError(t, mw.Close())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/gallery/upload", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp["data"].(map[string]interface{})["id"].(string)
}

func TestEditor_Gallery_UploadDeleteReorder(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Gallery Event")

	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/sections",
		map[string]interface{}{"gallery_enabled": true}, token).Code)

	id1 := uploadGalleryFile(t, server, token, eventID, "a.png")
	id2 := uploadGalleryFile(t, server, token, eventID, "b.png")

	wl := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/gallery", nil, token)
	require.Equal(t, http.StatusOK, wl.Code, wl.Body.String())
	var lresp map[string]interface{}
	require.NoError(t, json.Unmarshal(wl.Body.Bytes(), &lresp))
	photos := lresp["data"].([]interface{})
	require.Len(t, photos, 2)

	// Reorder
	reorder := map[string]interface{}{
		"photos": []map[string]interface{}{
			{"id": id2, "sort_order": 0},
			{"id": id1, "sort_order": 1},
		},
	}
	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/gallery/reorder", reorder, token).Code)

	// Verify new order
	wl2 := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/gallery", nil, token)
	var lresp2 map[string]interface{}
	require.NoError(t, json.Unmarshal(wl2.Body.Bytes(), &lresp2))
	photos2 := lresp2["data"].([]interface{})
	require.Equal(t, id2, photos2[0].(map[string]interface{})["id"].(string))

	// Delete
	wd := doRequest(t, server.Handler().(*chi.Mux), http.MethodDelete, "/api/v1/events/"+eventID+"/gallery/"+id1, nil, token)
	require.Equal(t, http.StatusOK, wd.Code, wd.Body.String())
	wl3 := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/gallery", nil, token)
	var lresp3 map[string]interface{}
	require.NoError(t, json.Unmarshal(wl3.Body.Bytes(), &lresp3))
	photos3 := lresp3["data"].([]interface{})
	require.Len(t, photos3, 1)

	// Public view renders gallery
	require.Equal(t, http.StatusOK, doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token).Code)
	slug := slugOf(t, server, token, eventID)
	_ = uploadGalleryFile(t, server, token, eventID, "c.png")
	wp := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/e/"+slug, nil, "")
	require.Equal(t, http.StatusOK, wp.Code, wp.Body.String())
	var pub map[string]interface{}
	require.NoError(t, json.Unmarshal(wp.Body.Bytes(), &pub))
	gallery := pub["data"].(map[string]interface{})["gallery"].([]interface{})
	require.Len(t, gallery, 2)
}

func TestEditor_Template_AssignAndList(t *testing.T) {
	server, db, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Template Event")

	// List templates available to free plan
	wl := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/templates", nil, token)
	require.Equal(t, http.StatusOK, wl.Code, wl.Body.String())
	var lresp map[string]interface{}
	require.NoError(t, json.Unmarshal(wl.Body.Bytes(), &lresp))
	templates := lresp["data"].([]interface{})
	require.NotEmpty(t, templates)

	// Assign a standard template
	tmpl := &model.Template{}
	require.NoError(t, db.Where("group_name = ?", "standard").First(tmpl).Error)
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/template",
		map[string]string{"template_id": tmpl.ID.String()}, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	wg := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	var gresp map[string]interface{}
	require.NoError(t, json.Unmarshal(wg.Body.Bytes(), &gresp))
	require.Equal(t, tmpl.ID.String(), gresp["data"].(map[string]interface{})["template_id"].(string))

	// Premium template forbidden for free plan
	premium := &model.Template{}
	require.NoError(t, db.Where("group_name = ?", "premium").First(premium).Error)
	wf := doRequest(t, server.Handler().(*chi.Mux), http.MethodPut, "/api/v1/events/"+eventID+"/template",
		map[string]string{"template_id": premium.ID.String()}, token)
	require.Equal(t, http.StatusForbidden, wf.Code, wf.Body.String())
}

func TestEditor_Music_PresetAssignAndUpload(t *testing.T) {
	server, db, userID, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Music Event")

	preset := &model.Music{Title: "Preset Song", Preset: "spotify:track:123", IsPreset: true}
	require.NoError(t, db.Create(preset).Error)

	wl := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/music/presets", nil, token)
	require.Equal(t, http.StatusOK, wl.Code, wl.Body.String())
	var lresp map[string]interface{}
	require.NoError(t, json.Unmarshal(wl.Body.Bytes(), &lresp))
	presets := lresp["data"].([]interface{})
	require.NotEmpty(t, presets)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/music/presets",
		map[string]string{"preset_id": preset.ID.String()}, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	wg := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, wg.Code, wg.Body.String())
	var gresp map[string]interface{}
	require.NoError(t, json.Unmarshal(wg.Body.Bytes(), &gresp))
	m := gresp["data"].(map[string]interface{})
	require.Equal(t, "Preset Song", m["title"].(string))
	require.Equal(t, "spotify:track:123", m["preset"].(string))
	require.Equal(t, true, m["is_preset"].(bool))

	// Upgrade to starter (has music.upload) for upload test
	starter := &model.Package{}
	require.NoError(t, db.Where("code = ?", "starter").First(starter).Error)
	require.NoError(t, db.Model(&model.Subscription{}).Where("user_id = ?", userID).Update("package_id", starter.ID).Error)

	eventID2 := createTestEvent(t, server, token, "Music Upload Event")
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	part, err := mw.CreateFormFile("file", "song.mp3")
	require.NoError(t, err)
	_, _ = part.Write([]byte("fake-audio"))
	_ = mw.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/"+eventID2+"/music/upload", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

func TestEditor_Music_NewlyCreatedEvent_ReturnsEmptyConfig(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "No Music Event")

	// A newly created event has no music configured. The endpoint must return
	// 200 with data:null, NOT an unexplained 404.
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, true, resp["success"])
	require.Nil(t, resp["data"], "expected empty music config (data: null)")
}

func TestEditor_Music_Remove_DisablesMusic(t *testing.T) {
	server, db, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Remove Music Event")

	// Assign a preset, then remove it.
	preset := &model.Music{Title: "Preset Song", Preset: "spotify:track:1", IsPreset: true}
	require.NoError(t, db.Create(preset).Error)
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/events/"+eventID+"/music/presets",
		map[string]string{"preset_id": preset.ID.String()}, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Verify it's selected.
	wg := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, wg.Code, wg.Body.String())

	// Remove the music.
	wd := doRequest(t, server.Handler().(*chi.Mux), http.MethodDelete, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, wd.Code, wd.Body.String())

	// After removal, GET returns empty config again (200 + null).
	wa := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, wa.Code, wa.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(wa.Body.Bytes(), &resp))
	require.Nil(t, resp["data"])

	// Removing again is a safe no-op.
	wd2 := doRequest(t, server.Handler().(*chi.Mux), http.MethodDelete, "/api/v1/events/"+eventID+"/music", nil, token)
	require.Equal(t, http.StatusOK, wd2.Code, wd2.Body.String())
}

func TestEditor_Music_Remove_Unauthorized(t *testing.T) {
	server, db, _, token := setupEditorTest(t)
	eventID := createTestEvent(t, server, token, "Owner Event")

	// Create a second user in the same DB and assign a subscription.
	otherUser := &model.User{Name: "Other", Email: "other-" + uuid.New().String() + "@test.com", PasswordHash: "x", Role: "user", Status: "active"}
	require.NoError(t, db.Create(otherUser).Error)
	freePkg := &model.Package{}
	require.NoError(t, db.Where("code = ?", "free").First(freePkg).Error)
	expiry := time.Now().Add(90 * 24 * time.Hour)
	require.NoError(t, db.Create(&model.Subscription{UserID: otherUser.ID, PackageID: freePkg.ID, Status: "active", StartAt: time.Now(), ExpiresAt: &expiry}).Error)
	otherToken := tokenForUser(t, otherUser.ID)

	// A different user may not remove music from an event they don't own.
	wd := doRequest(t, server.Handler().(*chi.Mux), http.MethodDelete, "/api/v1/events/"+eventID+"/music", nil, otherToken)
	require.Equal(t, http.StatusForbidden, wd.Code, wd.Body.String())
}

func TestEditor_Music_NotFoundEvent(t *testing.T) {
	server, _, _, token := setupEditorTest(t)
	missing := uuid.New().String()
	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/events/"+missing+"/music", nil, token)
	require.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}
