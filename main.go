package main

import (
	"fmt"

	"html/template"

	"log"

	"net/http"

	"os"

	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "./templ/" + p.Title + ".html"

	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./templ/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}

var validPath = regexp.MustCompile("^/(edit|save|blog)/([a-zA-Z0-9\055]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
	}
}

var templates = template.Must(template.ParseFiles("./static/edit.html", "./static/blog.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "blog", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/blog/"+title, http.StatusFound)
}

// TODO:
// Add default /blog/ view
// Render markdown
// Add a feed/log for blog posts
// Start thinking about things to write about

func main() {
	http.HandleFunc("/blog/", makeHandler(blogHandler))

	http.HandleFunc("/edit/", makeHandler(editHandler))

	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("Listening on: https://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
