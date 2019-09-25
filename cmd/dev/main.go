package main

import (
	"log"
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/handler"
	handler "github.com/neutral-af/backend/cmd/zeit"
)

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.URL.String(), r.Method, r.RemoteAddr)
		handler.Handler(w, r)
	})
	http.HandleFunc("/playground", gqlgenhandler.Playground("playground", "/graphql"))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
