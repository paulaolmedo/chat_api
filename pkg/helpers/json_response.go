package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondJSON translates an interface to json for response
func RespondJSON(w http.ResponseWriter, resp interface{}) {
	retJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(retJSON)
}

func JSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response) //puede ser datos de verdad o un error
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error()) //error parseando el mensaje
	}
}
