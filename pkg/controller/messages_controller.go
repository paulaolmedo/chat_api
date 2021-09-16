package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// SendMessage send a message from one user to another
func (config *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, JSONError)
		return
	}

	err = validateMessageData(message)
	if err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	timestamp := message.Timestamp
	message.Timestamp = timestamp.In(time.UTC)

	responseUser, err := config.Database.AddMessage(&message)
	if err != nil && (err.Error() == MissingSender || err.Error() == MissingRecipient) {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, UnknownError)
		return
	}

	helpers.JSONResponse(w, http.StatusCreated, responseUser)
}

// GetMessages get the messages from the logged user to a recipient
func (config *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// TODO: Retrieve list of Messages
	var filter models.MessageFilter
	var err error

	queryParams := r.URL.Query()

	errGettingParams := getQueryParams(queryParams, &filter)
	if errGettingParams != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, errGettingParams.Error())
		return
	}

	responseUser, err := config.Database.GetMessages(filter)
	if err != nil && err.Error() == MissingRecord {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, UnknownError)
		return
	}

	helpers.JSONResponse(w, http.StatusOK, responseUser)
}
