package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Blog struct {
	Posts      []Post
	Pagination Pagination
}

type Pagination struct {
	Page     int
	PrevPage int
	NextPage int
	Total    int
	Pages    []int
}

type Post struct {
	Id          string
	Title       string
	Description string
	Excerpt     string
	Slug        string
	Content     string
	CreatedAt   time.Time
}

type Router struct {
	chi.Router
	*Watcher
}

func New() *Router {
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

	w, err := NewWatcher(func(path string, b []byte) {
		var init bool

		err := filepath.Walk("web", func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, ".gohtml") {
				if !init {
					fmt.Println("Template updated.")
					tpl = template.New("")
					init = true
				}
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
	}, 10*time.Millisecond)
	if err != nil {
		panic(err)
	}

	w.Start()

	err = w.Watch("web")
	if err != nil {
		panic(err)
	}

	r := &Router{chi.NewRouter(), w}

	r.Use(middleware.Logger)

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

	renderPartial := func(partial string, data any) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err = tpl.ExecuteTemplate(w, partial, data)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	r.Get("/", renderApp("home"))
	r.Get("/blog", renderApp("blog"))
	r.Get("/projects", renderApp("projects"))

	r.Get("/partials/home", renderPartial("views/home", nil))
	r.Get("/partials/blog", func(w http.ResponseWriter, r *http.Request) {
		currentPage := 1
		page := r.URL.Query().Get("page")
		if page != "" {
			pag, err := strconv.ParseInt(page, 10, 64)
			if err != nil {
				panic(err)
			}

			if pag > 0 {
				currentPage = int(pag)
			}
		}

		var posts []Post
		for i := 0; i < 115; i++ {
			posts = append(posts, Post{
				Id:          fmt.Sprintf("uuid-%d", i),
				Title:       "My custom post " + fmt.Sprint(i),
				Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
				Excerpt:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
				Slug:        "my-custom-post-" + fmt.Sprint(i),
				Content:     "",
				CreatedAt:   time.Now(),
			})
		}

		perPage := 10

		blog := Blog{
			Posts: posts[currentPage*perPage : currentPage*perPage+10],
			Pagination: Pagination{
				Page:     currentPage,
				PrevPage: currentPage - 1,
				NextPage: currentPage + 1,
				Total:    12,
			},
		}

		blog.Pagination.UpdatePagination()

		b, _ := json.Marshal(blog.Pagination)
		fmt.Println(string(b))

		renderPartial("views/blog", blog)(w, r)
	})
	r.Get("/partials/projects", renderPartial("views/projects", nil))

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}

func (p *Pagination) UpdatePagination() {
	curr := p.Page

	first := max(1, curr-3)
	last := min(p.Total, curr+3)

	for i := first; i <= last; i++ {
		p.Pages = append(p.Pages, i)
	}

	for i := first - 1; i > max(last-7, 0); i-- {
		p.Pages = append([]int{i}, p.Pages...)
	}

	for i := last + 1; i <= min(7, p.Total); i++ {
		p.Pages = append(p.Pages, i)
	}

	if p.Pages[0] != 1 {
		p.Pages[0] = 1
	}

	if p.Pages[len(p.Pages)-1] != p.Total {
		p.Pages[len(p.Pages)-1] = p.Total
	}

	if len(p.Pages) > 1 {
		if p.Pages[1] != 2 {
			p.Pages[1] = 0
		}

		if p.Pages[len(p.Pages)-2] != p.Total-1 {
			p.Pages[len(p.Pages)-2] = 0
		}
	}
}
