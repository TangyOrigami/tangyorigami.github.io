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

// TODO:
// This needs to save to the ./templ/ folder

func (p *Page) save() error {
	filename := "templ/" + p.Title + ".html"

	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "templ/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")

}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/")

}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)

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
	http.Redirect(w, r, "/view/"+title, http.StatusFound)

}

var templates = template.Must(template.ParseFiles("./static/edit.html", "./static/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

// TODO:
// Can't handle this-syntax-here

func main() {
	http.HandleFunc("/view", makeHandler(viewHandler))

	http.HandleFunc("/edit", makeHandler(editHandler))

	http.HandleFunc("/save", makeHandler(saveHandler))

	// TODO:
	// Make this FileServer
	fserver := http.FileServer(http.Dir("."))

	log.Fatal(http.ListenAndServe(":8080", fserver))
	fmt.Printf("Listening on: https://localhost:8080")
}
