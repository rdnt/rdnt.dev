package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	chi.Router
}

func New() *Router {
	r := &Router{chi.NewRouter()}

	r.Use(middleware.Logger)

	//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello World!"))
	//})

	fs := http.FileServer(http.Dir("web"))

	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}
