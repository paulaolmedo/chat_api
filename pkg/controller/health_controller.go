package controller

import (
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// Check returns the health of the service and DB
func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	// TODO: Check service health. Feel free to add any check you consider necessary
	health := models.Health{Status: "ok"}
	helpers.JSONResponse(w, http.StatusOK, health)
}
