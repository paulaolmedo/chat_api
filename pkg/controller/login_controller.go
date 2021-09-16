package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/challenge/pkg/auth"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// Login authenticates a user and returns a bearer token
func (config *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, JSONError)
		return
	}

	if err := validateUserData(user); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	responseUser, err := config.Database.GetUser(user)
	if err != nil && err.Error() == MissingUser {
		helpers.JSONResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil { // de acá para abajo aún no hay tests
		helpers.JSONResponse(w, http.StatusInternalServerError, UnknownError)
		return
	}

	bearer, err := auth.GetBearerToken()
	if err != nil {
		errMsg := fmt.Errorf(TokenError, err)
		helpers.JSONResponse(w, http.StatusInternalServerError, errMsg)
		return
	}

	token := models.Login{Id: responseUser.UserID, Token: bearer.Token}

	helpers.JSONResponse(w, http.StatusCreated, token)
}
