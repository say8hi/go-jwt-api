package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/say8hi/go-jwt-api/internal/models"
	"github.com/stretchr/testify/assert"
)

var hash = sha256.Sum256([]byte("passwordtestuser"))
var authToken = hex.EncodeToString(hash[:])
var serverURL = "http://0.0.0.0:8081"

func TestCreateUserHandler_E2E(t *testing.T) {
	requestBody := models.CreateUserRequest{
		Username: "testuser",
	  Email: "email@email.com",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, serverURL+"/users/create", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var responseBody models.UserInDatabase
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
  expected := models.UserInDatabase{
    ID: "",
    Username: "testuser",
    Email: "email@email.com",
    RefreshHash: nil,
    IP: "",
}
  assert.NotEmpty(t, responseBody.ID, "ID should not be empty")
  expected.ID = responseBody.ID

  assert.NotEmpty(t, responseBody.IP, "ID should not be empty")
  expected.IP = responseBody.IP

  assert.Equal(t, expected, responseBody)
}

func TestCreateTokenPairs_E2E(t *testing.T) {
	requestBody := models.CreateUserRequest{
		Username: "testuser2",
	  Email: "email@email.com",
		IP: "192.168.2.3",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	createUserReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/create", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(createUserReq)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser models.UserInDatabase
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	assert.NoError(t, err)

	resp, err = http.Get(serverURL+"/users/get_tokens?user_id="+createdUser.ID)
	assert.NoError(t, err)

  type tokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
  }
  
  var responseBody tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

  assert.NotEmpty(t, responseBody.AccessToken)
  assert.NotEmpty(t, responseBody.RefreshToken)
}


func TestRefreshTokens_E2E(t *testing.T) {
	requestBody := models.CreateUserRequest{
		Username: "testuser3",
	  Email: "email@email.com",
		IP: "192.168.2.3",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	createUserReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/create", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(createUserReq)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser models.UserInDatabase
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	assert.NoError(t, err)

	resp, err = http.Get(serverURL+"/users/get_tokens?user_id="+createdUser.ID)
	assert.NoError(t, err)

  type tokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
  }
  
  var responseBody tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)


  requestRefreshBody := tokenResponse{
    AccessToken: responseBody.AccessToken,
    RefreshToken: responseBody.RefreshToken,
  }

	jsonData, err = json.Marshal(requestRefreshBody)

  req, err := http.NewRequest(http.MethodPost, serverURL+"/users/refresh", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	resp, err = client.Do(req)
	assert.NoError(t, err)

  var newResponseBody tokenResponse
  err = json.NewDecoder(resp.Body).Decode(&newResponseBody)
  assert.NoError(t, err)


  assert.NotEmpty(t, newResponseBody.AccessToken)
  assert.NotEmpty(t, newResponseBody.RefreshToken)
}
