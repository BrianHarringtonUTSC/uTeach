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
)

func getLogin(a *application.App, w http.ResponseWriter, r *http.Request) error {
	if _, ok := context.SessionUser(r); ok {
		return httperror.StatusError{http.StatusOK, errors.New("Already logged in")}
	}

	// TODO: replace code with CSRF token
	url := a.Config.OAuth2.AuthCodeURL("uteach-login") + "&connection=Username-Password-Authentication" // connection required by Auth0
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func loginUser(a *application.App, w http.ResponseWriter, r *http.Request, email, name string) error {
	um := models.NewUserModel(a.DB)
	user, err := um.FindOne(nil, squirrel.Eq{"users.email": email})
	if err == sql.ErrNoRows {
		// user not found so must be logging in for first time, add the user
		user = &models.User{Email: email, Name: name}
		err = um.Add(nil, user)
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
	errorParam := r.FormValue("error")
	if errorParam != "" {
		return httperror.StatusError{http.StatusUnauthorized,
			errors.Errorf("%s: %s", errorParam, r.FormValue("error_description"))}
	}

	// handle the exchange code to initiate a transport
	code := r.FormValue("code")
	if code == "" {
		return httperror.StatusError{http.StatusBadRequest, errors.New("Missing code.")}
	}

	tok, err := a.Config.OAuth2.Exchange(oauth2.NoContext, code)
	if err != nil {
		return httperror.StatusError{http.StatusUnauthorized, errors.New("Permission not given by user")}
	}

	// make get request to get user info using token
	client := a.Config.OAuth2.Client(oauth2.NoContext, tok)
	response, err := client.Get(a.Config.OAuth2UserInfoURL)
	if err != nil {
		return errors.Wrap(err, "get response error")
	}
	defer response.Body.Close()

	u := struct {
		Email string
		Name  string `json:"nickname"`
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
