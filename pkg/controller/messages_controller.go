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
	// TODO acá falta agregar quizás una struct adicional para que se adapte al json del swagger original
	// además, una validación que noté a último momento es que se corresponda el tipo de contenido que se está mandando, con los parámetros ingresados
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

	messageResponse, err := config.Database.AddMessage(&message)
	if err != nil && (err.Error() == MissingSender || err.Error() == MissingRecipient) {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, UnknownError)
		return
	}

	helpers.JSONResponse(w, http.StatusCreated, messageResponse)
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

	messages, err := config.Database.GetMessages(filter)
	if err != nil && err.Error() == MissingRecord {
		helpers.JSONResponse(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, UnknownError)
		return
	}

	// TODO y acá falta, al igual que en SendMessages, agregarle también una struct adicional para que se muestre como en el swagger original
	helpers.JSONResponse(w, http.StatusOK, messages)
}
