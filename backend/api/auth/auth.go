package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexedwards/scs/v2"
	verifier "github.com/okta/okta-jwt-verifier-golang"
)

type contextKey struct{}

var authCtxKey = &contextKey{}

type User struct {
	Email      string
	FamilyName string
	GivenName  string
}

func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(authCtxKey).(*User)
	return raw
}

func generateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

type Authenticator struct {
	ClientID       string
	ClientSecret   string
	Issuer         string
	RedirectURL    string
	SessionManager *scs.SessionManager
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm := a.SessionManager

		var user *User
		if email := sm.GetString(r.Context(), "email"); email != "" {
			user = &User{
				Email:      email,
				FamilyName: sm.GetString(r.Context(), "family_name"),
				GivenName:  sm.GetString(r.Context(), "given_name"),
			}
		}
		ctx := context.WithValue(r.Context(), authCtxKey, user)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	nonce, _ := generateNonce()
	state, _ := generateNonce()
	sm := a.SessionManager
	sm.Put(r.Context(), "auth-nonce", nonce)
	sm.Put(r.Context(), "auth-state", state)

	q := r.URL.Query()
	q.Add("client_id", a.ClientID)
	q.Add("nonce", nonce)
	q.Add("redirect_uri", a.RedirectURL)
	q.Add("response_mode", "query")
	q.Add("response_type", "code")
	q.Add("scope", "openid profile email")
	q.Add("state", state)

	redirectPath := a.Issuer + "/v1/authorize?" + q.Encode()
	http.Redirect(w, r, redirectPath, http.StatusFound)
}

func (a *Authenticator) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = a.SessionManager.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Authenticator) AuthCodeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "callback missing code", http.StatusInternalServerError)
		return
	}

	accessToken, idToken := a.getTokens(code, r)

	ctx := r.Context()
	sm := a.SessionManager
	state := sm.PopString(ctx, "auth-state")
	if state != r.URL.Query().Get("state") {
		http.Error(w, "callback state mismatch", http.StatusInternalServerError)
		return
	}

	accessTokenVerifier := verifier.JwtVerifier{
		Issuer: a.Issuer,
		ClaimsToValidate: map[string]string{
			"aud": "api://default",
			"cid": a.ClientID,
		},
	}
	result, err := accessTokenVerifier.New().VerifyAccessToken(accessToken)
	switch {
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case result == nil:
		http.Error(w, "access_token validation failed", http.StatusInternalServerError)
		return
	}

	idTokenVerifier := verifier.JwtVerifier{
		Issuer: a.Issuer,
		ClaimsToValidate: map[string]string{
			"aud":   a.ClientID,
			"nonce": sm.PopString(ctx, "auth-nonce"),
		},
	}
	result, err = idTokenVerifier.New().VerifyIdToken(idToken)
	switch {
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case result == nil:
		http.Error(w, "id_token validation failed", http.StatusInternalServerError)
		return
	}

	data := a.getProfileData(accessToken)
	sm.Put(r.Context(), "email", data["email"])
	sm.Put(r.Context(), "given_name", data["given_name"])
	sm.Put(r.Context(), "family_name", data["family_name"])
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Authenticator) getTokens(code string, r *http.Request) (accessToken, idToken string) {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(a.ClientID + ":" + a.ClientSecret),
	)

	q := r.URL.Query()
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", a.RedirectURL)
	url := a.Issuer + "/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte{}))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var tmp = struct {
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		AccessToken      string `json:"access_token,omitempty"`
		TokenType        string `json:"token_type,omitempty"`
		ExpiresIn        int    `json:"expires_in,omitempty"`
		Scope            string `json:"scope,omitempty"`
		IDToken          string `json:"id_token,omitempty"`
	}{}
	_ = json.Unmarshal(body, &tmp)
	return tmp.AccessToken, tmp.IDToken
}

func (a *Authenticator) getProfileData(accessToken string) map[string]string {
	url := a.Issuer + "/v1/userinfo"
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+accessToken)
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	data := map[string]string{}
	_ = json.Unmarshal(body, &data)
	return data
}
