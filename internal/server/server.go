package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Zucke/ContactManager/internal/handlers"
	"github.com/Zucke/ContactManager/pkg/authentication"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//Server is a struct that contains the param for the terver  like the port and the server of object
type Server struct {
	server *http.Server
	port   string
}

//Start put the server to listen
func (serv *Server) Start() {

	log.Printf("Escuchando en http://localhost:%s", serv.port)
	log.Fatal(serv.server.ListenAndServe())

}

func (serv *Server) getRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/register", handlers.NewUser)
	r.Post("/login", handlers.LoginUser)
	contacsGroup := r.Group(nil)
	contacsGroup.Use(authentication.ValidateMiddleware)
	contacsGroup.Post("/contacts/getall", handlers.GetAllContacts)
	contacsGroup.Post("/contacts/{name}", handlers.FDUContactByName)
	contacsGroup.Delete("/contacts/{name}", handlers.FDUContactByName)
	contacsGroup.Delete("/contacts/deleteall", handlers.DeleteAll)
	contacsGroup.Delete("/deleteme", handlers.DeleteUser)
	contacsGroup.Put("/contacts/{name}", handlers.FDUContactByName)
	contacsGroup.Post("/contacts/insertone", handlers.AddContact)

	return r

}

//New initialize the params for the server
func New(port string) *Server {
	serv := &Server{
		port: port,
	}

	r := serv.getRoutes()

	serv.server = &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return serv
}

//Close kill the server
func (serv *Server) Close(ctx context.Context) {
	log.Fatal(serv.server.Shutdown(ctx))
}
