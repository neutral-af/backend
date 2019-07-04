package main

import (
	"log"
	"net/http"

	handler "github.com/jasongwartz/carbon-offset-backend/cmd/zeit"
)

func main() {
	http.HandleFunc("/cmd/zeit/handler.go", handler.Handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
