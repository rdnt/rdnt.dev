package router

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	chi.Router
}

func New() *Router {
	r := &Router{chi.NewRouter()}

	r.Use(middleware.Logger)

	tpl := template.New("")
	err := filepath.Walk("web", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			_, err = tpl.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})
	if err != nil {
		panic(err)
	}

	type State struct {
		View string
	}

	renderApp := func(view string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err = tpl.ExecuteTemplate(w, "templates/main", State{View: view})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	renderPartial := func(partial string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err = tpl.ExecuteTemplate(w, partial, nil)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	r.Get("/", renderApp("home"))
	r.Get("/blog", renderApp("blog"))
	r.Get("/projects", renderApp("projects"))

	r.Get("/partials/home", renderPartial("views/home"))
	r.Get("/partials/blog", renderPartial("views/blog"))
	r.Get("/partials/projects", renderPartial("views/projects"))

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}
