package authguard

import (
	"net/http"

	jwttoken "ultimate.com/exercise/jwtToken"
	response "ultimate.com/exercise/responses"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := jwttoken.ValidateToken(r)
		if err != nil {
			response.RespondWithError(w, http.StatusBadRequest, "Invalid Token, Not authorized")
			return
		}
		endpoint(w, r)
	})
}
