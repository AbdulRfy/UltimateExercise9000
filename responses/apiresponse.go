package responses

import (
	"encoding/json"
	"net/http"

	models "ultimate.com/exercise/Models"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	updatedMessage := "Task Failed successfully: " + message
	RespondWithJSON(w, code, map[string]string{"error": updatedMessage})
}

func RespondWithMessage(w http.ResponseWriter, code int, message string) {

	defaultResponse := models.Response{}
	defaultResponse.Message = message
	RespondWithJSON(w, code, defaultResponse)
}

func RespondWithMessageAndToken(w http.ResponseWriter, code int, message string, token string) {

	defaultResponse := models.Response{}
	defaultResponse.Message = message
	defaultResponse.Token = token
	RespondWithJSON(w, code, defaultResponse)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
