package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(">", r.URL)
	}))
	appengine.Main()
}
