package main

import (
	"encoding/gob"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"strconv"
)

func handleSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := GetSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	RenderTemplate(w, r, "subjects.html", data)
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]

	topics, err := GetTopics(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Topics": topics}
	RenderTemplate(w, r, "topics.html", data)
}

func handleThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]

	threads, err := GetThreads(subjectName, topicName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userUpvotedThreadIDs := make(map[int]bool)
	user, ok := getSessionUser(r)
	if ok {
		userUpvotedThreadIDs, err = GetUserUpvotedThreadIDs(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data := map[string]interface{}{"Threads": threads, "UserUpvotedThreadIDs": userUpvotedThreadIDs}
	RenderTemplate(w, r, "threads.html", data)
}

func handleThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := GetThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Thread": thread}
	RenderTemplate(w, r, "thread.html", data)
}

func handleUpvote(w http.ResponseWriter, r *http.Request, fn func(string, int) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := getSessionUser(r)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	err = fn(user.Username, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func handleAddUpvote(w http.ResponseWriter, r *http.Request) {
	handleUpvote(w, r, AddUpVote)
}

func handleRemoveUpvote(w http.ResponseWriter, r *http.Request) {
	handleUpvote(w, r, RemoveUpvote)
}

func main() {
	LoadTemplates()

	err := InitDB()
	if err != nil {
		panic(err)
	}

	// allows user to be encoded so that it can be stored in a session
	gob.Register(&User{})

	authMiddleWare := alice.New(isAuth)

	router := mux.NewRouter()
	router.HandleFunc("/", handleSubjects)
	router.HandleFunc("/subjects", handleSubjects)
	router.HandleFunc("/topics/{subjectName}", handleTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", handleThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", handleThread)

	router.HandleFunc("/login/{username}", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	router.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(handleAddUpvote)).Methods("POST")
	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(handleRemoveUpvote)).Methods("DELETE")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
