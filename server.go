package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("tmpl/*.html"))

func handleSubjects(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "subjects.html", GetSubjects())
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	templates.ExecuteTemplate(w, "topics.html", GetTopics(subjectName))
}

func handleThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]
	templates.ExecuteTemplate(w, "threads.html", GetThreads(subjectName, topicName))
}

func handleThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]
	threadName := vars["threadName"]
	templates.ExecuteTemplate(w, "thread.html", GetThread(subjectName, topicName, threadName))
}

func main() {
	authMiddleWare := alice.New(isAuth)

	r := mux.NewRouter()

	r.HandleFunc("/", handleSubjects)
	r.HandleFunc("/subjects", handleSubjects)
	r.HandleFunc("/topics/{subjectName}", handleTopics)
	r.HandleFunc("/threads/{subjectName}/{topicName}", handleThreads)
	r.HandleFunc("/thread/{subjectName}/{topicName}/{threadName}", handleThread)

	r.HandleFunc("/login/{utorid}", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
