package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Zucke/ContactManager/internal/data"
	"github.com/Zucke/ContactManager/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetAllContacts return all contact for a user
func GetAllContacts(w http.ResponseWriter, r *http.Request) {

	contacts := data.NewUserContact()
	results, err := contacts.GetAll(r.Context())

	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, "Invalid Token")

	} else {
		render.Status(r, http.StatusFound)
		render.JSON(w, r, results)
	}

}

//AddContact add a contact to a especific user
func AddContact(w http.ResponseWriter, r *http.Request) {
	var contact data.Contact

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if contact.Name == "" || contact.Number == "" {
		response.HTTPError(w, r, http.StatusBadRequest, "Invalid information")
		return
	}

	contact.UserID = r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID)

	var userContact = data.NewUserContact()

	c, err := userContact.GetByName(r.Context(), contact.Name)
	if err == nil {
		response.HTTPError(w, r, http.StatusFound, "This contact already exist")
		render.JSON(w, r, c)
		return
	}

	c, err = userContact.GetByNumber(r.Context(), contact.Number)
	if err == nil {
		response.HTTPError(w, r, http.StatusFound, "This number already asociated to a contact name")
		render.JSON(w, r, c)
		return
	}

	err = userContact.InsertContact(contact)

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	} else {
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, contact)
	}

}

//FDUContactByName this find, delete or update a contact depent of the method used on request
func FDUContactByName(w http.ResponseWriter, r *http.Request) {
	var contact data.Contact
	var err error

	name := chi.URLParam(r, "name")
	userContact := data.NewUserContact()
	ctx := r.Context()

	switch r.Method {

	case http.MethodPost:
		contact, err = userContact.GetByName(ctx, name)

	case http.MethodDelete:
		contact, err = userContact.DeleteByName(ctx, name)

	case http.MethodPut:
		var updatedContact data.Contact
		err := json.NewDecoder(r.Body).Decode(&updatedContact)
		if err != nil {
			response.HTTPError(w, r, http.StatusFound, "error parsing new contact information")
			return
		}

		if updatedContact.Name == "" || updatedContact.Number == "" {
			response.HTTPError(w, r, http.StatusFound, "Invalid information")
			return
		}

		contact, err = userContact.UpdateContactByName(ctx, name, &updatedContact)

	default:
		return
	}

	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	render.Status(r, http.StatusFound)
	render.JSON(w, r, contact)

}

//DeleteAll delete all contact for a user
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	userContact := data.NewUserContact()
	err := userContact.DeleteAll(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	render.Status(r, http.StatusFound)

}
