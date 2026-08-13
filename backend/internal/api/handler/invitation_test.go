package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvitation_Create(t *testing.T) {
	setupAuthTestServer(t)

	resp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Invitation User",
		"email":    "invitation_user@example.com",
		"password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, resp.Code)

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "invitation_user@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":        "Wedding Andi dan Siti",
		"couple_name":  "Andi & Siti",
		"groom_name":   "Andi Pratama",
		"bride_name":   "Siti Rahayu",
		"wedding_date": "2026-08-15",
		"wedding_time": "10:00",
		"ceremony_venue": "Gedung Pernikahan",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)

	var eventResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &eventResp))
	require.True(t, eventResp["success"].(bool))
	data := eventResp["data"].(map[string]interface{})
	require.NotEmpty(t, data["id"])
	require.NotEmpty(t, data["slug"])
	require.Equal(t, "draft", data["status"])
}

func TestInvitation_Create_WithCustomSlug(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Slug User",
		"email":    "slug_user@example.com",
		"password": "securepassword123",
	}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "slug_user@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "My Wedding",
		"slug":       "custom-wedding-slug",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)

	var eventResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &eventResp))
	data := eventResp["data"].(map[string]interface{})
	require.Equal(t, "custom-wedding-slug", data["slug"])
}

func TestInvitation_Create_ReservedSlug(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "reserved_slug@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "reserved_slug@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Admin Wedding",
		"slug":       "admin",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	require.Equal(t, http.StatusCreated, w.Code)

	var eventResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &eventResp))
	data := eventResp["data"].(map[string]interface{})
	require.NotEqual(t, "admin", data["slug"])
}

func TestInvitation_GetByID(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "get_by_id@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "get_by_id@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Get Event Test",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	require.Equal(t, http.StatusOK, w.Code)

	var eventResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &eventResp))
	require.True(t, eventResp["success"].(bool))
}

func TestInvitation_Update(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "update_test@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "update_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Update Test",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodPut, "/api/v1/events/"+eventID, map[string]string{
		"title":            "Updated Wedding Title",
		"ceremony_venue":   "New Venue",
		"ceremony_address": "New Address",
	}, token)
	require.Equal(t, http.StatusOK, w.Code)

	var eventResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &eventResp))
	data := eventResp["data"].(map[string]interface{})
	require.Equal(t, "Updated Wedding Title", data["title"])
}

func TestInvitation_Delete(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "delete_test@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "delete_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Delete Test",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID, nil, token)
	require.Equal(t, http.StatusOK, w.Code)

	w2 := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	require.Equal(t, http.StatusNotFound, w2.Code)
}

func TestInvitation_Publish(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "publish_test@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "publish_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":        "Publish Test",
		"groom_name":   "Andi",
		"bride_name":   "Siti",
		"wedding_date": "2026-08-15",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)
	require.Equal(t, http.StatusOK, w.Code)

	var pubResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pubResp))
	require.True(t, pubResp["success"].(bool))
	data := pubResp["data"].(map[string]interface{})
	require.Equal(t, "published", data["status"])
	require.NotEmpty(t, data["published_at"])
	require.NotEmpty(t, data["public_url"])
}

func TestInvitation_Unpublish(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "unpublish_test@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "unpublish_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":        "Unpublish Test",
		"groom_name":   "Andi",
		"bride_name":   "Siti",
		"wedding_date": "2026-08-15",
	}, token)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/publish", nil, token)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/unpublish", nil, token)
	require.Equal(t, http.StatusOK, w.Code)

	getResp := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token)
	var eventData map[string]interface{}
	json.Unmarshal(getResp.Body.Bytes(), &eventData)
	data := eventData["data"].(map[string]interface{})
	require.Equal(t, "unpublished", data["status"])
}

func TestInvitation_OwnershipEnforced(t *testing.T) {
	setupAuthTestServer(t)

	resp1 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Owner User",
		"email":    "owner_user@example.com",
		"password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, resp1.Code)

	resp2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Other User",
		"email":    "other_user@example.com",
		"password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, resp2.Code)

	login1 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "owner_user@example.com",
		"password": "securepassword123",
	}, "")
	var login1Data map[string]interface{}
	json.Unmarshal(login1.Body.Bytes(), &login1Data)
	token1 := login1Data["data"].(map[string]interface{})["access_token"].(string)

	login2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "other_user@example.com",
		"password": "securepassword123",
	}, "")
	var login2Data map[string]interface{}
	json.Unmarshal(login2.Body.Bytes(), &login2Data)
	token2 := login2Data["data"].(map[string]interface{})["access_token"].(string)

	createResp := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Owner Event",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token1)
	var createData map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createData)
	eventID := createData["data"].(map[string]interface{})["id"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID, nil, token2)
	require.Equal(t, http.StatusForbidden, w.Code)

	w2 := doAuthRequest(t, http.MethodPut, "/api/v1/events/"+eventID, map[string]string{
		"title": "Hacked Title",
	}, token2)
	require.Equal(t, http.StatusForbidden, w2.Code)

	w3 := doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID, nil, token2)
	require.Equal(t, http.StatusForbidden, w3.Code)
}

func TestInvitation_List(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "list_test@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "list_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	for i := 0; i < 3; i++ {
		doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
			"title":      "List Event",
			"groom_name": "Andi",
			"bride_name": "Siti",
		}, token)
	}

	w := doAuthRequest(t, http.MethodGet, "/api/v1/events", nil, token)
	require.Equal(t, http.StatusOK, w.Code)

	var listResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &listResp))
	require.True(t, listResp["success"].(bool))
}

func TestInvitation_Publish_Validation(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "pub_validation@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "pub_validation@example.com",
		"password": "securepassword123",
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

func TestInvitation_SlugUniqueness(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "slug_unique@example.com",
		"password": "securepassword123",
		}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "slug_unique@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	resp1 := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Same Slug Test",
		"slug":       "unique-slug-test",
		"groom_name": "Andi",
		"bride_name": "Siti",
	}, token)
	require.Equal(t, http.StatusCreated, resp1.Code)

	resp2 := doAuthRequest(t, http.MethodPost, "/api/v1/events", map[string]string{
		"title":      "Same Slug Test 2",
		"slug":       "unique-slug-test",
		"groom_name": "Budi",
		"bride_name": "Ani",
	}, token)
	require.Equal(t, http.StatusCreated, resp2.Code)

	var data1, data2 map[string]interface{}
	json.Unmarshal(resp1.Body.Bytes(), &data1)
	json.Unmarshal(resp2.Body.Bytes(), &data2)
	require.NotEqual(t, data1["data"].(map[string]interface{})["slug"], data2["data"].(map[string]interface{})["slug"])
}
