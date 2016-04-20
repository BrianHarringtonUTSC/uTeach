package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/libtemplate"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

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

func getLogin(a *application.App, w http.ResponseWriter, r *http.Request) error {
	if _, ok := context.SessionUser(r); ok {
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
	user, err := u.GetUserByEmail(nil, email)

	if err == sql.ErrNoRows {
		// user not found so must be logging in for first time, add the user
		user, err = u.AddUser(nil, email, name)
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	us := session.NewUserSession(a.Store)
	err = us.SaveSessionUserID(w, r, user.ID)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func getOauth2Callback(a *application.App, w http.ResponseWriter, r *http.Request) error {
	googleConfig := getGoogleConfig(a, r)

	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := googleConfig.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		return httperror.StatusError{http.StatusUnauthorized, errors.New("Permission not given by user.")}
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

func getLogout(a *application.App, w http.ResponseWriter, r *http.Request) error {
	us := session.NewUserSession(a.Store)
	if err := us.Delete(w, r); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func getUser(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	email := strings.ToLower(vars["email"])

	um := models.NewUserModel(a.DB)
	user, err := um.GetUserByEmail(nil, email)
	if err != nil {
		return err
	}

	pm := models.NewPostModel(a.DB)
	createdPosts, err := pm.Find(nil, squirrel.Eq{"posts.creator_user_id": user.ID})
	if err != nil {
		return err
	}

	data := context.TemplateData(r)
	data["User"] = user
	data["CreatedPosts"] = createdPosts
	if err = addUserUpvotedPostIDsToData(r, pm, data); err != nil {
		return err
	}
	return libtemplate.Render(w, a.Templates, "user.html", data)
}
