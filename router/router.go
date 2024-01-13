package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
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
	CurrentPage        int
	TotalPages         int
	PerPage            int
	VisiblePages       []int
	First              int
	Last               int
	MaxVisible         int
	VisiblePagesMinus2 int
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
	r.Get("/partials/blog", func(writer http.ResponseWriter, request *http.Request) {
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

		currentPage := 4
		perPage := 10

		blog := Blog{
			Posts: posts[currentPage*perPage : currentPage*perPage+10],
			Pagination: Pagination{
				CurrentPage:        currentPage,
				TotalPages:         12,
				PerPage:            perPage,
				VisiblePages:       []int{},
				VisiblePagesMinus2: 0,
				First:              0,
				Last:               0,
				MaxVisible:         4, // divisible by 2
			},
		}

		blog.Pagination.UpdatePagination()

		b, _ := json.Marshal(blog.Pagination)
		fmt.Println(string(b))

		renderPartial("views/blog", blog)(writer, request)
	})
	r.Get("/partials/projects", renderPartial("views/projects", nil))

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}

func (p *Pagination) UpdatePagination() {
	p.First = int(math.Max(
		1,
		math.Min(
			float64(p.CurrentPage-p.MaxVisible/2),
			float64(p.TotalPages-p.MaxVisible),
		),
	))
	p.Last = int(math.Min(float64(p.First+p.MaxVisible), float64(p.TotalPages)))

	p.VisiblePages = []int{}
	for i := p.First; i <= p.Last; i++ {
		p.VisiblePages = append(p.VisiblePages, i)
	}
	if p.First > 1 {
		p.VisiblePages[0] = 1
	}
	if p.Last < p.TotalPages {
		p.VisiblePages[len(p.VisiblePages)-1] = p.TotalPages
	}

	p.VisiblePagesMinus2 = len(p.VisiblePages) - 2
}
