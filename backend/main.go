package main

//go:generate go run scripts/gqlgen.go

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/alexedwards/scs/v2"

	"github.com/sjansen/strongbox-elm/backend/api/auth"
	"github.com/sjansen/strongbox-elm/backend/api/graph"
)

const defaultPort = "8080"

var (
	clientID     = os.Getenv("STRONGBOX_OAUTH_CLIENT_ID")
	clientSecret = os.Getenv("STRONGBOX_OAUTH_CLIENT_SECRET")
	issuer       = os.Getenv("STRONGBOX_OAUTH_ISSUER")
	redirectURL  = os.Getenv("STRONGBOX_OAUTH_REDIRECT_URL")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	sm := scs.New()
	sm.Cookie.Name = "strongbox-session"
	sm.Cookie.Persist = true
	sm.IdleTimeout = 30 * time.Minute
	sm.Lifetime = 3 * time.Hour
	//sm.Cookie.Domain = "example.com"
	//sm.Cookie.HttpOnly = true
	//sm.Cookie.SameSite = http.SameSiteStrictMode
	//sm.Cookie.Secure = true
	if dynamostoreEndpoint != "" {
		store, err := newDynamoStore(dynamostoreEndpoint)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		sm.Store = store
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL playground", "/api/graphql"))
	mux.Handle("/api/graphql", handler.GraphQL(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})),
	)
	mux.HandleFunc("/session", getHandler)
	// auth
	auth := auth.Authenticator{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Issuer:         issuer,
		RedirectURL:    redirectURL,
		SessionManager: sm,
	}
	mux.HandleFunc("/login/oauth/callback", auth.AuthCodeCallbackHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/logout", auth.LogoutHandler)

	handler := sm.LoadAndSave(
		auth.Middleware(mux),
	)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, handler))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	user := auth.ForContext(r.Context())
	fmt.Fprintln(w, user.GivenName, user.FamilyName)
	fmt.Fprintln(w, user.Email)
}
