package controller

import (
	"encoding/json"
	"net/http"

	"github.com/challenge/pkg/auth"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// Login authenticates a user and returns a bearer token
func (config *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, "Error reading JSON data")
		return
	}

	if err := validateUserData(user); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	responseUser, err := config.Database.GetUser(user)
	if err != nil && err.Error() == "user not found" {
		helpers.JSONResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, "Error retrieving user information")
		return
	}

	bearer, err := auth.GetBearerToken()
	if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	token := models.Login{Id: responseUser.UserID, Token: bearer.Token}

	helpers.JSONResponse(w, http.StatusCreated, token)
}
