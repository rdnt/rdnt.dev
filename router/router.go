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

	r.Get("/", homePage)
	r.Get("/partials/home", homePagePartial)

	r.Get("/blog", blogPage)
	r.Get("/partials/blog", blogPagePartial)

	r.Get("/projects", projectsPage)
	r.Get("/partials/projects", projectsPagePartial)

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}
