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

	session := scs.NewSession()
	session.Cookie.Name = "strongbox-session"
	session.Cookie.Persist = true
	session.IdleTimeout = 30 * time.Minute
	session.Lifetime = 3 * time.Hour
	//session.Cookie.Domain = "example.com"
	//session.Cookie.HttpOnly = true
	//session.Cookie.SameSite = http.SameSiteStrictMode
	//session.Cookie.Secure = true

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL playground", "/api/graphql"))
	mux.Handle("/api/graphql", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))
	mux.HandleFunc("/session", getHandler)
	// auth
	auth := auth.Authenticator{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Issuer:       issuer,
		RedirectURL:  redirectURL,
		Session:      session,
	}
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/logout", auth.LogoutHandler)
	mux.HandleFunc("/authorization-code/callback", auth.AuthCodeCallbackHandler)

	handler := session.LoadAndSave(
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
