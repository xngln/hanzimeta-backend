package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/cors"
	"github.com/xngln/hanzimeta-backend/db"
	"github.com/xngln/hanzimeta-backend/graph"
	"github.com/xngln/hanzimeta-backend/graph/generated"
)

const defaultPort = "8080"

func main() {
	env := os.Getenv("HANZIMETA_ENV")
	if env != "PRODUCTION" {
		err := godotenv.Load()
		if err != nil {
			log.Panic(err)
		}
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// new relic setup
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Panic(err)
	}

	db.InitDB()

	mux := http.NewServeMux()

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),
	)

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle(newrelic.WrapHandle(newRelicApp, "/query", srv))

	var c *cors.Cors
	if env == "PRODUCTION" {
		c = cors.New(cors.Options{
			AllowedOrigins: []string{
				os.Getenv("UI_ORIGIN"),
			},
		})
	} else {
		c = cors.Default()
	}
	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
