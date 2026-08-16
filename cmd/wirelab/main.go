package main

import (
	"log"
	"net/http"
	"os"

	"wirelab.local/cabling/internal/articles"
)

func main() {
	data, err := articles.LoadYAML(articles.Fixture)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := articles.NewHandler(articles.NewService(data), http.FileServer(http.Dir("web")))
	log.Printf("wirelab listening on http://127.0.0.1:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
