package router

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func blogPage(w http.ResponseWriter, r *http.Request) {
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

	tpl, err := template.ParseFiles(
		"web/templates/partials/header.gohtml",
		"web/templates/partials/footer.gohtml",
		"web/templates/main.gohtml",
		"web/templates/pages/blog.gohtml",
	)
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(w, "main", blogPageData(currentPage))
	if err != nil {
		panic(err)
	}
}

func blogPagePartial(w http.ResponseWriter, r *http.Request) {
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

	tpl, err := template.ParseFiles(
		"web/templates/partial.gohtml",
		"web/templates/pages/blog.gohtml",
	)
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(w, "partial", blogPageData(currentPage))
	if err != nil {
		panic(err)
	}
}

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

type Post struct {
	Id          string
	Title       string
	Description string
	Excerpt     string
	Slug        string
	Content     string
	CreatedAt   time.Time
}

func blogPageData(currentPage int) Blog {
	perPage := 10
	totalPages := 100

	var posts []Post
	for i := 0; i < totalPages*perPage-perPage/2; i++ {
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

	blog := Blog{
		Posts: posts[currentPage*perPage : currentPage*perPage+10],
		Pagination: Pagination{
			Page:     currentPage,
			PrevPage: currentPage - 1,
			NextPage: currentPage + 1,
			Total:    totalPages,
		},
	}

	blog.Pagination.UpdatePagination()

	return blog
}
