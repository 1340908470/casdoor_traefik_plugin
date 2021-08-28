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
	next http.Handler
	name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UIDDemo{
		next: next,
		name: name,
	}, nil
}

func (r *UIDDemo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello world"))
}
