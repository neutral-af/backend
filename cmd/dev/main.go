package main

import (
	"log"
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/handler"
	handler "github.com/jasongwartz/carbon-offset-backend/cmd/zeit"
)

func main() {
	http.HandleFunc("/graphql", handler.Handler)
	http.HandleFunc("/playground", gqlgenhandler.Playground("playground", "/graphql"))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
