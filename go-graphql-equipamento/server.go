package main

import (
	"go-graphql-equipamento/graph"
	"go-graphql-equipamento/graph/generated"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

var serverloger = log.New(os.Stdout, "ServerGate.......{*} ", log.LstdFlags)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	serverloger.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	serverloger.Fatal(http.ListenAndServe(":"+port, nil))
}
