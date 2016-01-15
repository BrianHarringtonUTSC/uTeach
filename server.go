package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("tmpl/*.html"))

func handleSubjects(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "subjects.html", GetSubjects())
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	// TODO: refactor repeated handle functoins into a generic function
	vars := mux.Vars(r)

	subjectID, err := strconv.ParseInt(vars["subjectID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templates.ExecuteTemplate(w, "topics.html", GetTopics(subjectID))
}

func handleThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	topicID, err := strconv.ParseInt(vars["topicID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templates.ExecuteTemplate(w, "threads.html", GetThreads(topicID))
}

func handleThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templates.ExecuteTemplate(w, "thread.html", GetThread(threadID))
}

func main() {
	authMiddleWare := alice.New(isAuth)

	r := mux.NewRouter()

	r.HandleFunc("/", handleSubjects)
	r.HandleFunc("/topics/{subjectID}", handleTopics)
	r.HandleFunc("/threads/{topicID}", handleThreads)
	r.HandleFunc("/thread/{threadID}", handleThread)

	r.HandleFunc("/login/{utorid}", handleLogin)
	r.HandleFunc("/logout", handleLogout)
	r.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
