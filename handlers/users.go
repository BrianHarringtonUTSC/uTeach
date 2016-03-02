package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// getGoogleConfig gets an oauth2 Config for doing authentication with Google.
func getGoogleConfig(r *http.Request) *oauth2.Config {
	app := application.GetFromContext(r)

	googleConfig := &oauth2.Config{
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	googleConfig.RedirectURL = app.Config.GoogleRedirectURL
	googleConfig.ClientID = app.Config.GoogleClientID
	googleConfig.ClientSecret = app.Config.GoogleClientSecret

	return googleConfig
}

// GetLogin makes a request to Google Oauth2 authenticator.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	if _, ok := app.Store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	googleConfig := getGoogleConfig(r)
	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := googleConfig.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	http.Redirect(w, r, url, http.StatusFound)
}

// GetOauth2Callback responds to callbacks from Google Oauth2 authenticator.
func GetOauth2Callback(w http.ResponseWriter, r *http.Request) {
	googleConfig := getGoogleConfig(r)

	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := googleConfig.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// make get request to get user info using token
	client := googleConfig.Client(oauth2.NoContext, tok)
	response, err := client.Get(googleUserInfoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// parse response into generic map
	m := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new user session
	app := application.GetFromContext(r)
	email := m["email"].(string)

	err = app.Store.NewUserSession(w, r, email, app.DB) // username is email for now
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout logs the user out.
func Logout(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	if err := app.Store.DeleteUserSession(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// GetUser renders user info.
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)
	userCreatedThreads, err := t.GetThreadsByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Username": username, "UserCreatedThreads": userCreatedThreads}
	renderTemplate(w, r, "user.html", data)
}
