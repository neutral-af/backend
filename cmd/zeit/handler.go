package handler

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/jasongwartz/carbon-offset-backend/lib/schema"
	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

var graphQLHandler http.HandlerFunc

func init() {
	graphQLHandler = handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &schema.Resolver{}}))
	fmt.Println("Registered graphQL handler")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	graphQLHandler(w, r)
}
