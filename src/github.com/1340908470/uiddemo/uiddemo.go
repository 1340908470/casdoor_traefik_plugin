package uiddemo

import (
	"context"
	"net/http"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type UIDDemo struct {
	next http.Handler // next middleware's handler
	name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UIDDemo{
		next: next,
		name: name,
	}, nil
}

func (r *UIDDemo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.next != nil {
		r.next.ServeHTTP(rw, req)
	} else {
		r.next = http.DefaultServeMux
	}
}

/*

# Forward auth middleware

## 1. root handler
- Modify request
- Adjust if we're acting as forward auth middleware
- Pass to mux (router handler)

> https://github.com/thomseddon/traefik-forward-auth/blob/master/internal/server.go :57

## 2. router handler

*/
