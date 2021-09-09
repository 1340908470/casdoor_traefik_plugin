package casdoorauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-session/session"
)

type Config struct {
	// when login using casdoor successfully, the casdoor would jump to this URL
	RedirectURI string
	// get client id from casdoor - application - edit page
	ClientID string
	// host where casdoor service deployed
	ServiceHost string
}

func CreateConfig() *Config {
	return &Config{}
}

type CasdoorAuth struct {
	next        http.Handler // next middleware's handler
	name        string
	callbackURL *url.URL
	clientID    string
	serviceHost string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	callbackURL, _ := url.Parse(config.RedirectURI)

	return &CasdoorAuth{
		next:        next,
		name:        name,
		callbackURL: callbackURL,
		clientID:    config.ClientID,
		serviceHost: config.ServiceHost,
	}, nil
}

func (r *CasdoorAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	store, _ := session.Start(context.Background(), rw, req)

	if req.URL.Path == r.callbackURL.Path {
		code := req.URL.Query().Get("code")
		state := req.URL.Query().Get("state")

		if code != "" && state != "" {
			token, err := auth.GetOAuthToken(code, state)
			if err != nil {
				panic(err)
			}

			claims, err := auth.ParseJwtToken(token.AccessToken)
			if err != nil {
				panic(err)
			}

			claims.AccessToken = token.AccessToken

			store.Set("casdoor_claims", claims)
			store.Set("casdoor_claims", token)
			_ = store.Save()
		}

	} else {
		claims, isok := store.Get("casdoor_claims")
		if !isok {
			loginURLStr := fmt.Sprintf("%v/login/oauth/authorize?client_id=%v&response_type=code&redirect_uri=%v&scope=read&state=casdoor", r.serviceHost, r.clientID, r.callbackURL)
			http.Redirect(rw, req, loginURLStr, http.StatusFound)
			return
		} else {
			rw.Header().Add("casdoor_claims", fmt.Sprintf("%v", claims))
		}
	}

	if r.next != nil {
		r.next.ServeHTTP(rw, req)
	} else {
		r.next = http.DefaultServeMux
	}
}
