package main

//go:generate go run scripts/gqlgen.go

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"

	"github.com/sjansen/strongbox-elm/backend/api/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/api/graphql"))
	http.Handle("/api/graphql", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}
