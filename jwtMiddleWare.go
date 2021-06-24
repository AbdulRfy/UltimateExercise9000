package main

import (
	"net/http"
)

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := validateToken(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Token, Not authorized")
			return
		}
		endpoint(w, r)
	})
}
