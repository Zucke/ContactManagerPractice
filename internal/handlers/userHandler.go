package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zucke/ContactManager/internal/data"
	"github.com/Zucke/ContactManager/pkg/authentication"
	"github.com/Zucke/ContactManager/pkg/response"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const expiryIn = 2 * time.Minute

var tokenAuth *jwtauth.JWTAuth

//LoginUser login a user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var result, requestUser authentication.User
	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, "Failed to parse user")
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	usersData := data.NewUserData()
	result, err = usersData.VarifyUserByNickname(ctx, requestUser.Nickname)

	if err != nil {
		response.HTTPError(w, r, http.StatusForbidden, "Bad information")
		return
	}

	if !requestUser.ComparePassword(result.Password) {
		response.HTTPError(w, r, http.StatusForbidden, "Bad information")
		return
	}
	var token string
	result.Password = ""
	token, err = authentication.GenerateJWT(result)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	result.Password = ""
	render.JSON(w, r, render.M{
		"user":  result,
		"token": token,
	})

}

//NewUser create a new user
func NewUser(w http.ResponseWriter, r *http.Request) {
	var newUser *authentication.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, "Failed to parse user")
		return
	}

	if newUser.Password == "" || newUser.Nickname == "" {
		response.HTTPError(w, r, http.StatusBadRequest, "invalid information")
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	usersData := data.NewUserData()
	err = usersData.NewUser(ctx, newUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"user": newUser.Nickname})

}

//DeleteUser delete the corrend logged user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userData := data.NewUserData()
	userData.DeleteUser(r.Context())
}
