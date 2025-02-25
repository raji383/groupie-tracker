package main

import (
	"log"
	"net/http"
	web "Web/Backend"
	// "strings"
	// "os"
	// "fmt"
)
func main() {
	http.HandleFunc("/filter", web.FilterHandler)
	http.HandleFunc("/search", web.SearchHandler)

	// css := http.FileServer(http.Dir("./Frontend/css"))
	http.HandleFunc("/css/", web.Test)
	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/artist", web.ArtistHandler)
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

