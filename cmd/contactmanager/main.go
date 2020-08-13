package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Zucke/ContactManager/internal/server"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8000"
	}
	serv := server.New(port)

	go serv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Finalizando Conexion...")

	ctx := context.Background()
	serv.Close(ctx)

}
