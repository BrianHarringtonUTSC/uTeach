package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"

	"github.com/umairidris/uTeach/application"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// credentials should be obtained from the Google Developer Console (https://console.developers.google.com).
var conf = &oauth2.Config{
	ClientID:     os.Getenv("UTEACH_GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("UTEACH_GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/oauth2callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

// GetLogin makes a request to Google Oauth2 authenticator.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	if _, ok := app.Store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	if conf.ClientID == "" || conf.ClientSecret == "" {
		http.Error(w, "Google Oauth2 Client ID and/or secret not set.", http.StatusInternalServerError)
		return
	}

	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := conf.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	http.Redirect(w, r, url, http.StatusFound)
}

// GetOauth2Callback responds to callbacks from Google Oauth2 authenticator.
func GetOauth2Callback(w http.ResponseWriter, r *http.Request) {
	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := conf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// make get request to get user info using token
	client := conf.Client(oauth2.NoContext, tok)
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

	// create new user session using info from user info
	app := application.Get(r)
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
	app := application.Get(r)
	if err := app.Store.DeleteUserSession(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
