package test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(test *testing.T) {
	user := InitUser("testUser", "12345678")

	client := resty.New()
	response, err := client.R().
		SetBody(user).
		EnableTrace().
		Post(testServer.URL + "/users")

	require.NoError(test, err)
	assert.Equal(test, http.StatusCreated, response.StatusCode())
}

func TestCreateUserThatAlreadyExists(test *testing.T) {
	user := InitUser("testUser", "12345678")

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
