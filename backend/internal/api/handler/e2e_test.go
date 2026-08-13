package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var e2eCounter = 0

func e2eEmail(prefix string) string {
	e2eCounter++
	return fmt.Sprintf("%s_%d@example.com", prefix, e2eCounter)
}

func TestE2E_RegisterLoginCreateEvent(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("e2e")
	regResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "E2E User", "email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, regResp.Code)

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.Code)

	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)
	require.NotEmpty(t, token)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "E2E Wedding", "groom_name": "Groom", "bride_name": "Bride", "wedding_date": "2026-08-15",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)
}

func TestE2E_UnauthorizedAccess(t *testing.T) {
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events", nil, "")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	w2 := doAuthRequest(t, http.MethodGet, "/api/v1/events", nil, "Bearer invalidtoken")
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestE2E_AdminForbiddenForUser(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("user")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Normal User", "email": email, "password": "securepassword123",
	}, "")
	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/admin/analytics", nil, token)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_DuplicateRegistration(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("dup")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User One", "email": email, "password": "securepassword123",
	}, "")

	dupResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User Two", "email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusConflict, dupResp.Code)
}

func TestE2E_InvalidLogin(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("invalid")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Valid User", "email": email, "password": "securepassword123",
	}, "")

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "wrongpassword",
	}, "")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	w2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": e2eEmail("nonexistent"), "password": "password123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestE2E_EventNotFound(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("notfound")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User", "email": email, "password": "securepassword123",
	}, "")
	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/00000000-0000-0000-0000-000000000000", nil, token)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestE2E_IDORProtection(t *testing.T) {
	setupAuthTestServer(t)

	email1 := e2eEmail("user1")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User One", "email": email1, "password": "securepassword123",
	}, "")
	loginResp1 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email1, "password": "securepassword123",
	}, "")
	var loginData1 map[string]interface{}
	json.Unmarshal(loginResp1.Body.Bytes(), &loginData1)
	token1 := loginData1["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "Private Event", "groom_name": "G", "bride_name": "B", "wedding_date": "2026-08-15",
	}, token1)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	email2 := e2eEmail("user2")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User Two", "email": email2, "password": "securepassword123",
	}, "")
	loginResp2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email2, "password": "securepassword123",
	}, "")
	var loginData2 map[string]interface{}
	json.Unmarshal(loginResp2.Body.Bytes(), &loginData2)
	token2 := loginData2["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token2)
	require.Equal(t, http.StatusForbidden, w.Code)

	w2 := doAuthRequest(t, http.MethodPut, "/api/v1/events/"+eventID, map[string]string{
		"title": "Hacked",
	}, token2)
	require.Equal(t, http.StatusForbidden, w2.Code)

	w3 := doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID, nil, token2)
	require.Equal(t, http.StatusForbidden, w3.Code)
}

func TestE2E_ValidationError(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("validation")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User", "email": email, "password": "securepassword123",
	}, "")
	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "",
	}, token)
	require.Equal(t, http.StatusBadRequest, w.Code)

	w2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email,
	}, token)
	require.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestE2E_PublishValidation(t *testing.T) {
	setupAuthTestServer(t)

	email := e2eEmail("publish")
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "User", "email": email, "password": "securepassword123",
	}, "")
	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title": "Incomplete Event",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
