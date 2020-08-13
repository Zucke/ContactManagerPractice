package handlers

import (
	"net/http"

	"github.com/Zucke/ContactManager/pkg/response"
	"github.com/dgrijalva/jwt-go"
)

//IsValidToken common function to validate and send error in case is not a valit token or a erro has occurred while parsing
func IsValidToken(token jwt.Token, err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
		return false
	}

	if !token.Valid {
		response.HTTPError(w, r, http.StatusUnauthorized, "Invalid Token")
		return false
	}

	return true
}
