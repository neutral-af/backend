package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	gqlgen_handler "github.com/99designs/gqlgen/handler"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/neutral-af/backend/lib/config"
	generated "github.com/neutral-af/backend/lib/graphql-generated"
	"github.com/neutral-af/backend/lib/resolvers"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/gqlerror"
)

var graphQLHandler http.HandlerFunc

func init() {
	beeline.Init(beeline.Config{
		WriteKey:    config.C.HoneycombAPIKey,
		Dataset:     "carbonara-backend",
		ServiceName: "carbonara-go",
	})
	// defer beeline.Close()

	graphQLHandler = hnynethttp.WrapHandlerFunc(
		gqlgen_handler.GraphQL(
			generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}),
			gqlgen_handler.ErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
				log.Printf("ERROR: %s", e)
				return graphql.DefaultErrorPresenter(ctx, e)
			}),
		),
	)
	log.Print("Registered graphQL handler")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var origins []string

	if config.C.Environment == config.EnvironmentProd {
		origins = []string{"https://neutral.af", "https://*.neutral-af.now.sh"}
	} else {
		origins = []string{"*"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins: origins,
	})

	c.Handler(http.HandlerFunc(graphQLHandler)).ServeHTTP(w, r)
}
