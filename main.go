package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/content", apiContentHandler) // Handle htmx request for dynamic content

	log.Print("Listening on localhost:3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func apiContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.Header.Get("Hx-Request") == "true" {
		// Handle htmx request for dynamic content
		content := "This is dynamic content from the backend!"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"content": content})
		return
	}

	// Handle non-htmx requests here if needed
	http.NotFound(w, r)
}

