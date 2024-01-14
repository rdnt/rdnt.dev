package router

import (
	"html/template"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(
		"web/templates/partials/header.gohtml",
		"web/templates/partials/footer.gohtml",
		"web/templates/main.gohtml",
		"web/templates/pages/home.gohtml",
	)
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(w, "main", nil)
	if err != nil {
		panic(err)
	}
}

func homePagePartial(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(
		"web/templates/partial.gohtml",
		"web/templates/pages/home.gohtml",
	)
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(w, "partial", nil)
	if err != nil {
		panic(err)
	}
}
