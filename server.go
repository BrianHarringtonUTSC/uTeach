package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Subject struct {
	ID   int64
	Name string
}

type Topic struct {
	ID        int64
	Name      string
	SubjectId int64
}

type Resource struct {
	ID             int64
	Name           string
	Content        string
	ThreadID       int64
	PostedByUserID int64
}

var templates = template.Must(template.ParseFiles("tmpl/subjects.html", "tmpl/topics.html"))

func handleSubjects(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "subjects.html", GetSubjects())
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectID, err := strconv.ParseInt(vars["subjectID"], 10, 64)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	templates.ExecuteTemplate(w, "topics.html", GetTopics(subjectID))
}

func main() {
	authMiddleWare := alice.New(isAuth)

	r := mux.NewRouter()

	r.HandleFunc("/", handleSubjects)
	r.HandleFunc("/topics/{subjectID}", handleTopics)

	r.HandleFunc("/login/{utorid}", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
