package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Title string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/api/content", apiContentHandler)

	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Your Website Title",
	}

	tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <header>
        <h1>{{.Title}}</h1>
    </header>
    <main>
        <section>
            <h2>Welcome to Your Website</h2>
            <p id="content" hx-get="/api/content">Loading content...</p>
            <!-- Add more content as needed -->
        </section>
    </main>
    <footer>
        <p>&copy; 2024 Your Website. All rights reserved.</p>
    </footer>
    <script src="https://unpkg.com/htmx.org@1.7.2/dist/htmx.min.js"></script>
    <script src="/static/scripts.js"></script>
</body>
</html>
`)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func apiContentHandler(w http.ResponseWriter, r *http.Request) {
	content := "This is dynamic content from the backend!"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": content})
}

