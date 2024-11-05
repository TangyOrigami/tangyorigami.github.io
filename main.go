package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"regexp"
	"strings"

	"github.com/go-mail/mail"

	"github.com/bmizerany/pat"
	"github.com/joho/godotenv"
)

/*
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

*/

func main() {
	godotenv.Load()
	mux := pat.New()

	fs := http.FileServer(http.Dir("assets"))
	mux.Get("/assets/", http.StripPrefix("/assets/", fs))
	mux.Get("/", http.HandlerFunc(home))
	mux.Post("/", http.HandlerFunc(send))
	mux.Get("/confirmation", http.HandlerFunc(confirmation))

	//http.HandleFunc("/blog/", makeHandler(blogHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Print("Listening...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Print(os.Getenv("THIS"))
	render(w, "static/home.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	// Validate form
	msg := &Message{
		Email:   r.PostFormValue("email"),
		Content: r.PostFormValue("content"),
	}

	if msg.Validate() == true {
		// Send contact form message in email
		if err := msg.DeliverAsync(); err != nil {
			log.Print(err)
			http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
			return
		}

		// Redirect to confirmation
		http.Redirect(w, r, "confirmation", http.StatusSeeOther)
	}
	if msg.Validate() == false {
		render(w, "static/home.html", msg)
		http.Redirect(w, r, "confirmation", http.StatusSeeOther)
	}

}

func confirmation(w http.ResponseWriter, r *http.Request) {
	// Add GIPHY integration, eventually

	render(w, "static/confirmation.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}

// EMAIL
var rxEmail = regexp.MustCompile(".+@.+\\..+")

type Message struct {
	Email   string
	Content string
	Errors  map[string]string
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(msg.Email))
	if match == false {
		msg.Errors["Email"] = "Please enter a valid email address."
	}

	if strings.TrimSpace(msg.Content) == "" {
		msg.Errors["Content"] = "Please enter a message"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) DeliverAsync() error {
	username := os.Getenv("MAILTRAP_UN")
	password := os.Getenv("MAILTRAP_PW")

	email := mail.NewMessage()
	email.SetHeader("To", "carlos@csaenz.dev")
	email.SetHeader("From", "server@csaenz.dev")
	email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "New message via Contact Form")
	email.SetBody("text/plain", msg.Content)
	email.AddAlternative("text/html",
		`
		<!DOCTYPE html>
		<html>
		<head>
		<style>
		body {
		font-family: Arial, sans-serif;
		color: #333333;
		background-color: #f9f9f9;
		margin: 0;
		padding: 0;
		}
		.email-container {
		width: 100%;
		max-width: 600px;
		margin: 20px auto;
		background-color: #ffffff;
		border: 1px solid #dddddd;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		padding: 20px;
		}
		h1 {
		color: #0073e6;
		font-size: 24px;
		text-align: center;
		}
		p {
		font-size: 16px;
		line-height: 1.5;
		}
		.content {
		padding: 10px;
		background-color: #f1f1f1;
		border-radius: 5px;
		margin: 15px 0;
		font-style: italic;
		}
		.footer {
		text-align: center;
		font-size: 14px;
		color: #888888;
		margin-top: 20px;
		}
		</style>
		</head>
		<body>
		<div class="email-container">
		<h1>New Message via Contact Form</h1>
		<p><strong>Email:</strong> `+msg.Email+`</p>
		<div class="content">"`+msg.Content+`"</div>
		<p class="footer">Made with &#128140;<br>Mailtrap</p>
		</div>
		</body>
		</html>
		`)

	return mail.NewDialer("live.smtp.mailtrap.io", 587, username, password).DialAndSend(email)
}
