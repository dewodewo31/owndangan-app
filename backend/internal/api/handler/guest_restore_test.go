package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestGuest(t *testing.T, token, eventID, name string) string {
	cr := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests", map[string]string{
		"name": name,
	}, token)
	require.Equal(t, http.StatusCreated, cr.Code, cr.Body.String())
	var cd map[string]interface{}
	require.NoError(t, json.Unmarshal(cr.Body.Bytes(), &cd))
	return cd["data"].(map[string]interface{})["id"].(string)
}

func listGuestIDs(t *testing.T, token, eventID, path string) []string {
	w := doAuthRequest(t, http.MethodGet, "/api/v1/events/"+eventID+"/guests"+path, nil, token)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	arr, _ := data["data"].([]interface{})
	ids := make([]string, 0, len(arr))
	for _, g := range arr {
		ids = append(ids, g.(map[string]interface{})["id"].(string))
	}
	return ids
}

func containsID(ids []string, id string) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

func TestGuest_SoftDelete_MovesToTrash(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	gid := createTestGuest(t, token, eventID, "Trash")

	del := doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token)
	require.Equal(t, http.StatusOK, del.Code)

	active := listGuestIDs(t, token, eventID, "")
	require.False(t, containsID(active, gid), "deleted guest must not appear in active list")
	deleted := listGuestIDs(t, token, eventID, "/deleted")
	require.True(t, containsID(deleted, gid), "deleted guest must appear in trash")
}

func TestGuest_Restore(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	gid := createTestGuest(t, token, eventID, "Restore")

	doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token)
	res := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/"+gid+"/restore", nil, token)
	require.Equal(t, http.StatusOK, res.Code)

	active := listGuestIDs(t, token, eventID, "")
	require.True(t, containsID(active, gid), "restored guest must reappear in active list")
	deleted := listGuestIDs(t, token, eventID, "/deleted")
	require.False(t, containsID(deleted, gid), "restored guest must leave trash")
}

func TestGuest_RestoreActive_Idempotent(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	gid := createTestGuest(t, token, eventID, "Active")

	res := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/"+gid+"/restore", nil, token)
	require.Equal(t, http.StatusOK, res.Code)
	active := listGuestIDs(t, token, eventID, "")
	require.True(t, containsID(active, gid), "active guest stays in active list after restore")
}

func registerUser(t *testing.T) string {
	testCounter++
	email := fmt.Sprintf("intruder_%d@example.com", testCounter)
	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Intruder", "email": email, "password": "securepassword123",
	}, "")
	lr := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusOK, lr.Code, lr.Body.String())
	var ld map[string]interface{}
	require.NoError(t, json.Unmarshal(lr.Body.Bytes(), &ld))
	return ld["data"].(map[string]interface{})["access_token"].(string)
}

func TestGuest_RestoreOtherUser_Forbidden(t *testing.T) {
	setupAuthTestServer(t)
	token1, eventID := createTestUserAndEvent(t)
	gid := createTestGuest(t, token1, eventID, "Victim")

	token2 := registerUser(t)
	res := doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/"+gid+"/restore", nil, token2)
	require.Equal(t, http.StatusForbidden, res.Code)
}

func TestGuest_DeleteRestoreDelete_Stable(t *testing.T) {
	setupAuthTestServer(t)
	token, eventID := createTestUserAndEvent(t)
	gid := createTestGuest(t, token, eventID, "Cycle")

	require.Equal(t, http.StatusOK, doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token).Code)
	require.Equal(t, http.StatusOK, doAuthRequest(t, http.MethodPost, "/api/v1/events/"+eventID+"/guests/"+gid+"/restore", nil, token).Code)
	require.Equal(t, http.StatusOK, doAuthRequest(t, http.MethodDelete, "/api/v1/events/"+eventID+"/guests/"+gid, nil, token).Code)
	deleted := listGuestIDs(t, token, eventID, "/deleted")
	require.True(t, containsID(deleted, gid), "guest must be trashable again after restore")
}
