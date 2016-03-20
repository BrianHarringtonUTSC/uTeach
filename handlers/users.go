package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"strings"

	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// getGoogleConfig gets an oauth2 Config for doing authentication with Google.
func getGoogleConfig(r *http.Request) *oauth2.Config {

	googleConfig := &oauth2.Config{
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	config := context.Config(r)
	googleConfig.RedirectURL = config.GoogleRedirectURL
	googleConfig.ClientID = config.GoogleClientID
	googleConfig.ClientSecret = config.GoogleClientSecret

	return googleConfig
}

// GetLogin makes a request to Google Oauth2 authenticator.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := getSessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	googleConfig := getGoogleConfig(r)
	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := googleConfig.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	http.Redirect(w, r, url, http.StatusFound)
}

func loginUser(w http.ResponseWriter, r *http.Request, email, name string) {
	u := models.NewUserModel(context.DB(r))
	user, err := u.GetUserByEmail(email)
	// sign up if user is logging in for first time
	if err == sql.ErrNoRows {
		user, err = u.Signup(email, name)
	}
	if err != nil {
		handleError(w, err)
		return
	}

	usm := session.NewUserSessionManager(context.CookieStore(r))
	err = usm.New(w, r, user)
	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// GetOauth2Callback responds to callbacks from Google Oauth2 authenticator.
func GetOauth2Callback(w http.ResponseWriter, r *http.Request) {
	googleConfig := getGoogleConfig(r)

	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := googleConfig.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		handleError(w, err)
		return
	}

	// make get request to get user info using token
	client := googleConfig.Client(oauth2.NoContext, tok)
	response, err := client.Get(googleUserInfoURL)
	if err != nil {
		handleError(w, err)
		return
	}
	defer response.Body.Close()

	// parse response into generic map
	m := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		handleError(w, err)
		return
	}

	email := m["email"].(string)
	name := m["name"].(string)
	loginUser(w, r, email, name)
}

// Logout logs the user out.
func Logout(w http.ResponseWriter, r *http.Request) {
	usm := session.NewUserSessionManager(context.CookieStore(r))
	if err := usm.Delete(w, r); err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// GetUser renders user info.
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := strings.ToLower(vars["email"])

	tm := getThreadModel(r)
	createdThreads, err := tm.GetThreadsByEmail(email)
	if err != nil {
		handleError(w, err)
		return
	}

	data := map[string]interface{}{"Email": email, "CreatedThreads": createdThreads}
	renderTemplate(w, r, "user.html", data)
}
