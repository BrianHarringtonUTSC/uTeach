package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"strings"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// getGoogleConfig gets an oauth2 Config for doing authentication with Google.
func getGoogleConfig(a *application.App, r *http.Request) *oauth2.Config {

	googleConfig := &oauth2.Config{
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	googleConfig.RedirectURL = a.Config.GoogleRedirectURL
	googleConfig.ClientID = a.Config.GoogleClientID
	googleConfig.ClientSecret = a.Config.GoogleClientSecret

	return googleConfig
}

// GetLogin makes a request to Google Oauth2 authenticator.
func GetLogin(a *application.App, w http.ResponseWriter, r *http.Request) error {
	usm := session.NewUserSessionManager(a.CookieStore)
	if _, ok := usm.SessionUser(r); ok {
		return errors.New("Already logged in")
	}

	googleConfig := getGoogleConfig(a, r)
	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := googleConfig.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func loginUser(a *application.App, w http.ResponseWriter, r *http.Request, email, name string) error {
	u := models.NewUserModel(a.DB)
	user, err := u.GetUserByEmail(email)

	// sign up if user is logging in for first time
	if err == sql.ErrNoRows {
		user, err = u.Signup(email, name)
	}
	if err != nil {
		return err
	}

	usm := session.NewUserSessionManager(a.CookieStore)
	err = usm.New(w, r, user)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

// GetOauth2Callback responds to callbacks from Google Oauth2 authenticator.
func GetOauth2Callback(a *application.App, w http.ResponseWriter, r *http.Request) error {
	googleConfig := getGoogleConfig(a, r)

	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := googleConfig.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		return err
	}

	// make get request to get user info using token
	client := googleConfig.Client(oauth2.NoContext, tok)
	response, err := client.Get(googleUserInfoURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// parse response into generic map
	m := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		return err
	}

	email := m["email"].(string)
	name := m["name"].(string)
	return loginUser(a, w, r, email, name)
}

// Logout logs the user out.
func GetLogout(a *application.App, w http.ResponseWriter, r *http.Request) error {
	usm := session.NewUserSessionManager(a.CookieStore)
	if err := usm.Delete(w, r); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

// GetUser renders user info.
func GetUser(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	email := strings.ToLower(vars["email"])

	um := models.NewUserModel(a.DB)
	user, err := um.GetUserByEmail(email)
	if err != nil {
		return err
	}

	tm := models.NewThreadModel(a.DB)
	createdThreads, err := tm.GetThreadsByUser(user)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Email": email, "CreatedThreads": createdThreads}
	return renderTemplate(a, w, r, "user.html", data)
}
