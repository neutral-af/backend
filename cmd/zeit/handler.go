package handler

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	generated "github.com/jasongwartz/carbon-offset-backend/lib/graphql-generated"
	"github.com/jasongwartz/carbon-offset-backend/lib/resolvers"
)

var graphQLHandler http.HandlerFunc

func init() {
	graphQLHandler = handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))
	fmt.Println("Registered graphQL handler")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	graphQLHandler(w, r)
}
