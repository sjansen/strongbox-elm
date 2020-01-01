package main

//go:generate go run scripts/gqlgen.go

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sjansen/strongbox-elm/backend/server"
)

const defaultPort = "8080"

func main() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	handler, err := server.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	handler = server.RequestLogger(handler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, handler))
}
