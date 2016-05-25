package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/context"
	"github.com/BrianHarringtonUTSC/uTeach/httperror"
	"github.com/BrianHarringtonUTSC/uTeach/libtemplate"
	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/BrianHarringtonUTSC/uTeach/session"
	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
		return httperror.StatusError{http.StatusOK, errors.New("Already logged in")}
	}

	googleConfig := getGoogleConfig(a, r)
	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := googleConfig.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func loginUser(a *application.App, w http.ResponseWriter, r *http.Request, email, name string) error {
	um := models.NewUserModel(a.DB)
	user, err := um.FindOne(nil, squirrel.Eq{"users.email": email})
	if err == sql.ErrNoRows {
		// user not found so must be logging in for first time, add the user
		user, err = um.AddUser(nil, email, name)
	}
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "login error")
	}

	us := session.NewUserSession(a.Store)
	err = us.SaveSessionUserID(w, r, user.ID)
	if err != nil {
		return errors.Wrap(err, "save session user error")
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
		return errors.Wrap(err, "get response error")
	}
	defer response.Body.Close()

	u := struct {
		Email string
		Name  string
	}{}

	err = json.NewDecoder(response.Body).Decode(&u)
	if err != nil {
		return errors.Wrap(err, "json decode error")
	}
	return loginUser(a, w, r, u.Email, u.Name)
}

func getLogout(a *application.App, w http.ResponseWriter, r *http.Request) error {
	us := session.NewUserSession(a.Store)
	if err := us.Delete(w, r); err != nil {
		return errors.Wrap(err, "user session delete error")
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func getUser(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	email := strings.ToLower(vars["email"])

	um := models.NewUserModel(a.DB)
	user, err := um.FindOne(nil, squirrel.Eq{"users.email": email})
	if err != nil {
		return errors.Wrap(err, "find one error")
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
		return errors.Wrap(err, "add upvoted post ids to data error")
	}
	err = libtemplate.Render(w, a.Templates, "user.html", data)
	return errors.Wrap(err, "render template error")
}
