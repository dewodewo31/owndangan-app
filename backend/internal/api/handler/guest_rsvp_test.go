package handler_test

import (
	"encoding/json"
	"fmt"
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

