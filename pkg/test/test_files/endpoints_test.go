package test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(test *testing.T) {
	user := InitUser("testUser1", "12345678")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusCreated, response.StatusCode())
}

func TestCreateUserThatAlreadyExists(test *testing.T) {
	user := InitUser("testUser1", "12345678")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusConflict, response.StatusCode())
}

func TestCreateUserWithCorruptedData(test *testing.T) {
	client := resty.New()
	response, err := client.R().
		SetBody(nil).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestCreateUserWithEmptyUsername(test *testing.T) {
	user := InitUser("", "12345678")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestCreateUserWithEmptyPassword(test *testing.T) {
	user := InitUser("testUser", "")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestLoginWithNonExistentUser(test *testing.T) {
	user := InitUser("nonExistentUser", "12345678")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/login")

	require.NoError(test, err)
	assert.Equal(test, http.StatusNotFound, response.StatusCode())
}

func TestLoginWithCorruptedData(test *testing.T) {
	client := resty.New()
	response, err := client.R().
		SetBody(nil).
		EnableTrace().
		Post(testServer.URL + "/login")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestCreateMessage(test *testing.T) {
	user := InitMessage(1, 1)

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusCreated, response.StatusCode())
}

func TestMessagesWithCorruptedData(test *testing.T) {
	client := resty.New()
	response, err := client.R().
		SetBody(nil).
		EnableTrace().
		Post(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}
func TestCreateMessageWithInvalidData(test *testing.T) {
	user := InitMessage(0, 1)

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestCreateMessageWithNonExistentUsers(test *testing.T) {
	user := InitMessage(11, 1)

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusConflict, response.StatusCode())
}

func TestSearchMessage(test *testing.T) {
	client := resty.New()
	response, err := client.R().EnableTrace().
		SetQueryParams(map[string]string{"recipient": "1", "start": "1", "limit": "1"}).
		Get(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusOK, response.StatusCode())
}

func TestSearchMessageOfNonExistentRecipient(test *testing.T) {
	client := resty.New()
	response, err := client.R().EnableTrace().
		SetQueryParams(map[string]string{"recipient": "123", "start": "1"}).
		Get(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusConflict, response.StatusCode())
}

func TestSearchMessageWithInvalidRecipient(test *testing.T) {
	client := resty.New()
	response, err := client.R().EnableTrace().
		SetQueryParams(map[string]string{"recipient": "-1", "start": "1"}).
		Get(testServer.URL + "/messages")

	require.NoError(test, err)
	assert.Equal(test, http.StatusBadRequest, response.StatusCode())
}

func TestHealthCheck(test *testing.T) {
	client := resty.New()
	response, err := client.R().
		EnableTrace().
		Post(testServer.URL + "/check")

	require.NoError(test, err)
	assert.Equal(test, http.StatusOK, response.StatusCode())
}
