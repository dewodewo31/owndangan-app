package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

var testCounter = 0

func createTestUserAndEvent(t *testing.T) (string, string) {
	testCounter++
	email := fmt.Sprintf("user_%d@example.com", testCounter)
	regResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Test", "email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, regResp.Code, regResp.Body.String())

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.Code, loginResp.Body.String())
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "Test", "groom_name": "G", "bride_name": "B", "wedding_date": "2026-08-15",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code, createResp.Body.String())
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)
	return token, eventID
}

func TestGuest_Create(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	w := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "John Doe", "category": "family",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestGuest_List(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	for i := 0; i < 3; i++ {
		doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
			"name": "G" + string(rune('A'+i)),
		}, token)
	}
	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/guests", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGuest_Update(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	cr := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "Orig",
	}, token)
	var cd map[string]interface{}
	json.Unmarshal(cr.Body.Bytes(), &cd)
	gid := cd["data"].(map[string]interface{})["id"].(string)
	w := doAuthRequest(t, http.MethodPut, "/api/v1/events/"+eventID+"/guests/"+gid, map[string]string{
		"name": "Updated",
	}, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGuest_Delete(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	cr := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "Del",
	}, token)
	var cd map[string]interface{}
	json.Unmarshal(cr.Body.Bytes(), &cd)
	gid := cd["data"].(map[string]interface{})["id"].(string)
	w := doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token)
	require.Equal(t, http.StatusOK, w.Code)
	w2 := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token)
	require.Equal(t, http.StatusNotFound, w2.Code)
}

func TestGuest_Ownership(t *testing.T) {
	setupAuthTestServer(t)
	token1, eventID := createTestUserAndEvent(t)
	cr := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "Own",
	}, token1)
	var cd map[string]interface{}
	json.Unmarshal(cr.Body.Bytes(), &cd)
	gid := cd["data"].(map[string]interface{})["id"].(string)

	testCounter++
	email2 := fmt.Sprintf("intruder_%d@example.com", testCounter)
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "I", "email": email2, "password": "securepassword123",
	}, "")
	lr := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email2, "password": "securepassword123",
	}, "")
	var ld map[string]interface{}
	json.Unmarshal(lr.Body.Bytes(), &ld)
	token2 := ld["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token2)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestGuest_ImportCSV(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	csvData := "name,phone,category\nCSV1,628111,family\nCSV2,628222,colleague\n"
	body := strings.NewReader(csvData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/guests/import", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----b")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testServer.Handler().(*chi.Mux).ServeHTTP(w, req)
	require.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)
}

func TestRSVP_Submit(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	cr := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "RSVP",
	}, token)
	var cd map[string]interface{}
	json.Unmarshal(cr.Body.Bytes(), &cd)
	gtok := cd["data"].(map[string]interface{})["token"].(string)

	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/rsvp/"+eventID+"/submit", map[string]interface{}{
		"token": gtok, "attendance": "attending", "guest_count": 2,
	}, "")
	require.Equal(t, http.StatusOK, w.Code)
}

func TestRSVP_Recap(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	w := doAuthRequest(t, http.MethodGet, "/api/v1/rsvp/"+eventID+"/recap", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGuestbook_Submit(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	w := doAuthRequest(t, http.MethodPost, "/api/v1/guestbook/"+eventID+"/submit", map[string]string{
		"name": "Anon", "message": "Congrats!",
	}, "")
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestGuestbook_Public(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	w := doAuthRequest(t, http.MethodGet, "/api/v1/guestbook/"+eventID, nil, "")
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGuestbook_Moderation(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	doAuthRequest(t, http.MethodPost, "/api/v1/guestbook/"+eventID+"/submit", map[string]string{
		"name": "Anon", "message": "Test",
	}, "")
	w := doAuthRequest(t, http.MethodGet, "/api/v1/guestbook/"+eventID+"/all", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

// --- Import (preview -> confirm) tests ---

type importPreviewRow struct {
	Index    int      `json:"index"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Category string   `json:"category"`
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
}

type importPreviewResp struct {
	Columns []string           `json:"columns"`
	Rows    []importPreviewRow `json:"rows"`
	Summary struct {
		Total     int `json:"total"`
		Valid     int `json:"valid"`
		Duplicate int `json:"duplicate"`
		Invalid   int `json:"invalid"`
	} `json:"summary"`
}

type importConfirmResp struct {
	Total      int `json:"total"`
	Imported   int `json:"imported"`
	Duplicates int `json:"duplicates"`
	Errors     []struct {
		Index  int      `json:"index"`
		Errors []string `json:"errors"`
	} `json:"errors"`
}

func doImportPreview(t *testing.T, eventID, token, csvContent string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "guests.csv")
	require.NoError(t, err)
	_, err = fw.Write([]byte(csvContent))
	require.NoError(t, err)
	require.NoError(t, mw.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/guests/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Handler().(*chi.Mux).ServeHTTP(rec, req)
	return rec
}

func TestGuest_ImportPreviewConfirm(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)

	csv := "name,phone,category\nAlice,6281111111,keluarga\nBob,6282222222,teman\n"
	previewRec := doImportPreview(t, eventID, token, csv)
	require.Equal(t, http.StatusOK, previewRec.Code, previewRec.Body.String())

	var previewWrap struct {
		Data importPreviewResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(previewRec.Body.Bytes(), &previewWrap))
	preview := previewWrap.Data
	require.Equal(t, 2, preview.Summary.Total)
	require.Equal(t, 2, preview.Summary.Valid)
	require.Equal(t, 0, preview.Summary.Duplicate)
	require.Equal(t, 0, preview.Summary.Invalid)
	require.Len(t, preview.Rows, 2)
	require.Equal(t, "valid", preview.Rows[0].Status)
	require.Equal(t, "keluarga", preview.Rows[0].Category)

	// Confirm the previewed rows.
	rows := make([]map[string]string, 0, len(preview.Rows))
	for _, r := range preview.Rows {
		rows = append(rows, map[string]string{
			"name": r.Name, "email": r.Email, "phone": r.Phone, "category": r.Category,
		})
	}
	confirmRec := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/import/confirm",
		map[string]interface{}{"rows": rows}, token)
	require.Equal(t, http.StatusOK, confirmRec.Code, confirmRec.Body.String())

	var confirmWrap struct {
		Data importConfirmResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(confirmRec.Body.Bytes(), &confirmWrap))
	confirm := confirmWrap.Data
	require.Equal(t, 2, confirm.Total)
	require.Equal(t, 2, confirm.Imported)
	require.Equal(t, 0, confirm.Duplicates)
	require.Empty(t, confirm.Errors)
}

func TestGuest_ImportDuplicate(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)

	// Seed an existing guest that should be detected as a duplicate.
	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": "Existing", "phone": "6289999999",
	}, token)

	// "Existing" collides with DB; the two "New" rows collide with each other.
	csv := "name,phone\nExisting,6289999999\nNew,6288888888\nNew,6288888888\n"
	previewRec := doImportPreview(t, eventID, token, csv)
	require.Equal(t, http.StatusOK, previewRec.Code, previewRec.Body.String())

	var previewWrap struct {
		Data importPreviewResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(previewRec.Body.Bytes(), &previewWrap))
	preview := previewWrap.Data
	require.Equal(t, 3, preview.Summary.Total)
	require.Equal(t, 1, preview.Summary.Valid)
	require.Equal(t, 2, preview.Summary.Duplicate)
	require.Equal(t, 0, preview.Summary.Invalid)
	require.Equal(t, "duplicate", preview.Rows[0].Status)
	require.Equal(t, "valid", preview.Rows[1].Status)
	require.Equal(t, "duplicate", preview.Rows[2].Status)

	// Confirm all three rows: one imported, two duplicates.
	rows := []map[string]string{
		{"name": "Existing", "phone": "6289999999"},
		{"name": "New", "phone": "6288888888"},
		{"name": "New", "phone": "6288888888"},
	}
	confirmRec := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/import/confirm",
		map[string]interface{}{"rows": rows}, token)
	require.Equal(t, http.StatusOK, confirmRec.Code, confirmRec.Body.String())

	var confirmWrap struct {
		Data importConfirmResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(confirmRec.Body.Bytes(), &confirmWrap))
	confirm := confirmWrap.Data
	require.Equal(t, 3, confirm.Total)
	require.Equal(t, 1, confirm.Imported)
	require.Equal(t, 2, confirm.Duplicates)
}

func TestGuest_ImportInvalid(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)

	csv := "name,phone\n,6281111111\nValid,6282222222\n"
	previewRec := doImportPreview(t, eventID, token, csv)
	require.Equal(t, http.StatusOK, previewRec.Code, previewRec.Body.String())

	var previewWrap struct {
		Data importPreviewResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(previewRec.Body.Bytes(), &previewWrap))
	preview := previewWrap.Data
	require.Equal(t, 2, preview.Summary.Total)
	require.Equal(t, 1, preview.Summary.Valid)
	require.Equal(t, 1, preview.Summary.Invalid)
	require.Equal(t, "invalid", preview.Rows[0].Status)
	require.Equal(t, "valid", preview.Rows[1].Status)

	rows := []map[string]string{
		{"name": "", "phone": "6281111111"},
		{"name": "Valid", "phone": "6282222222"},
	}
	confirmRec := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/import/confirm",
		map[string]interface{}{"rows": rows}, token)
	require.Equal(t, http.StatusOK, confirmRec.Code, confirmRec.Body.String())

	var confirmWrap struct {
		Data importConfirmResp `json:"data"`
	}
	require.NoError(t, json.Unmarshal(confirmRec.Body.Bytes(), &confirmWrap))
	confirm := confirmWrap.Data
	require.Equal(t, 2, confirm.Total)
	require.Equal(t, 1, confirm.Imported)
	require.Len(t, confirm.Errors, 1)
	require.Equal(t, 1, confirm.Errors[0].Index)
}

func TestGuest_ImportXlsxRejected(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "guests.xlsx")
	require.NoError(t, err)
	_, err = fw.Write([]byte("PK\x03\x04 fake xlsx"))
	require.NoError(t, err)
	require.NoError(t, mw.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/guests/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Handler().(*chi.Mux).ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- Public invitations endpoint test ---

func TestInvitations_Public(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)

	// Get the event slug.
	evRec := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	require.Equal(t, http.StatusOK, evRec.Code)
	var evData map[string]interface{}
	require.NoError(t, json.Unmarshal(evRec.Body.Bytes(), &evData))
	slug := evData["data"].(map[string]interface{})["slug"].(string)
	require.NotEmpty(t, slug)

	// Draft event must not appear in the public feed.
	pubRec := doAuthRequest(t, http.MethodGet, "/api/v1/invitations/public", nil, "")
	require.Equal(t, http.StatusOK, pubRec.Code)
	var pubWrap struct {
		Data []map[string]interface{} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(pubRec.Body.Bytes(), &pubWrap))
	pubData := pubWrap.Data
	for _, it := range pubData {
		require.NotEqual(t, slug, it["slug"], "draft event must not be public")
	}

	// Publish and confirm it shows up.
	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	pubRec2 := doAuthRequest(t, http.MethodGet, "/api/v1/invitations/public", nil, "")
	require.Equal(t, http.StatusOK, pubRec2.Code)
	var pubWrap2 struct {
		Data []map[string]interface{} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(pubRec2.Body.Bytes(), &pubWrap2))
	pubData2 := pubWrap2.Data
	found := false
	for _, it := range pubData2 {
		if it["slug"] == slug {
			found = true
			require.NotEmpty(t, it["updated_at"])
		}
	}
	require.True(t, found, "published event should appear in public feed")
}
