package main

import (
	"log"
	"net/http"

	"github.com/napster-mayu/service-inspector/route"
)

func main() {
	serverPort := ":8080"

	r := route.GetRoutes()

	err := make(chan error)
	log.Println("successgully started server on port: ", serverPort)
	err <- http.ListenAndServe(serverPort, r)

	log.Printf("received error on servcer port %v: %v", serverPort, <-err)
	log.Fatalf("shutting down server on: %v", serverPort)
}
