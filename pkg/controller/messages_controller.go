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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, "Error reading JSON data")
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
	if err != nil && (err.Error() == "sender does not exist" || err.Error() == "recipient does not exist") {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, "Error retrieving user information")
		return
	}

	helpers.JSONResponse(w, http.StatusCreated, responseUser)
}

// GetMessages get the messages from the logged user to a recipient
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// TODO: Retrieve list of Messages
	helpers.RespondJSON(w, []*models.Message{{}})
}
