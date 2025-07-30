package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a file server to serve files from the "static" directory.
	// The path given to http.Dir is relative to the project directory.
	fileServer := http.FileServer(http.Dir("."))
	static := http.FileServer(http.Dir("./static"))
	assets := http.FileServer(http.Dir("./assets"))
	blog := http.FileServer(http.Dir("./blog"))

	// Register the file server as the handler for all URL paths.
	// "/" matches all request paths.
	http.Handle("/", fileServer)
	http.Handle("/static", static)
	http.Handle("/assets", assets)
	http.Handle("/blog", blog)

	log.Println("Listening on :8080...")
	// Start the server and listen on port 8080.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
