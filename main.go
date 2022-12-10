package main

import (
	"fmt"
	"net/http"
)

// custom mux
type CustomMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &c.ServeMux

	for _, next := range c.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

func (c *CustomMux) RegisterMiddleware(middleware func(next http.Handler) http.Handler) {
	c.middlewares = append(c.middlewares, middleware)
}

// handler
func ShowIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// middleware
func AllowOnlyGET(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		isValid := (username == "benzion") && (password == "yehezkel")
		if !ok || !isValid {
			http.Error(w, "Wrong username or password!", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := new(CustomMux)
	mux.HandleFunc("/", ShowIndex)

	var handler http.Handler = mux
	mux.RegisterMiddleware(AllowOnlyGET)
	mux.RegisterMiddleware(Authenticate)

	server := new(http.Server)
	server.Addr = ":3000"
	server.Handler = handler

	fmt.Println("server running on port 3000")
	server.ListenAndServe()
}
