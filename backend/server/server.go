package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/alexedwards/scs/v2"
	"github.com/sjansen/strongbox-elm/backend/api/auth"
	"github.com/sjansen/strongbox-elm/backend/api/graph"
)

var (
	clientID     = os.Getenv("STRONGBOX_OAUTH_CLIENT_ID")
	clientSecret = os.Getenv("STRONGBOX_OAUTH_CLIENT_SECRET")
	issuer       = os.Getenv("STRONGBOX_OAUTH_ISSUER")
	redirectURL  = os.Getenv("STRONGBOX_OAUTH_REDIRECT_URL")
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("x-amzn-trace-id")
		log.Printf("%-4s %q %s\n", r.Method, r.URL.Path, id)
		next.ServeHTTP(w, r)
	})
}

func New() (http.Handler, error) {
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
			return nil, err
		}
		sm.Store = store
	}

	auth := auth.Authenticator{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Issuer:         issuer,
		RedirectURL:    redirectURL,
		SessionManager: sm,
	}

	graphql := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL Playground", "/api/graphql"))
	mux.Handle("/api/graphql", handler.GraphQL(graphql))
	mux.HandleFunc("/api/login", auth.LoginHandler)
	mux.HandleFunc("/api/login/oauth/callback", auth.AuthCodeCallbackHandler)
	mux.HandleFunc("/api/logout", auth.LogoutHandler)

	handler := sm.LoadAndSave(
		auth.Middleware(mux),
	)

	return handler, nil
}
