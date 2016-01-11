package main

import (
	"html/template"
	"log"
	"net/http"
)

type Topic struct {
	Name string
}

var templates = template.Must(template.ParseFiles("tmpl/topics.html"))

var topics = []Topic{Topic{"Python"}, Topic{"Java"}} // TODO: switch to a database

func handleTopics(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "topics.html", topics)
}

func main() {
	http.HandleFunc("/", handleTopics)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
