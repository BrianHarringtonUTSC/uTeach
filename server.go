package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Topic struct {
	Name string
}

// views
var templates = template.Must(template.ParseFiles("tmpl/topics.html"))

// temp hard coded topics
var topics = []Topic{Topic{"Python"}, Topic{"Java"}} // TODO: switch to a database

func handleTopics(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "topics.html", topics)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleTopics)

	r.HandleFunc("/login/{utorid}", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/check", handleCheck)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
