package handler

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/jasongwartz/carbon-offset-backend/lib/config"
	generated "github.com/jasongwartz/carbon-offset-backend/lib/graphql-generated"
	"github.com/jasongwartz/carbon-offset-backend/lib/resolvers"
	"github.com/rs/cors"
)

var graphQLHandler http.HandlerFunc

func init() {
	beeline.Init(beeline.Config{
		WriteKey:    config.C.HoneycombAPIKey,
		Dataset:     "carbonara-backend",
		ServiceName: "carbonara-go",
	})
	// defer beeline.Close()

	graphQLHandler = hnynethttp.WrapHandlerFunc(handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}})))
	fmt.Println("Registered graphQL handler")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var origins []string

	if config.C.Environment == "prod" {
		origins = []string{"https://frontend.jasongwartz.now.sh"}
	} else {
		origins = []string{"http://localhost:8080"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins: origins,
	})

	c.Handler(http.HandlerFunc(graphQLHandler)).ServeHTTP(w, r)
}
