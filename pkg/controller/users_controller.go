package controller

import (
	"encoding/json"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// CreateUser creates a new user given an username and a password
func (config *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, "Error reading JSON data")
		return
	}

	err = validateUserData(user)
	if err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	responseUser, err := config.Database.CreateUser(user)
	if err != nil && err.Error() == "user already exists" {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(w, http.StatusCreated, responseUser)
}
