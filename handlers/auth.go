package handlers

import (
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"net/http"
	"io/ioutil"

	"github.com/umairidris/uTeach/application"
)

// credentials should be obtained from the Google Developer Console (https://console.developers.google.com).
var	conf = &oauth2.Config {
	    ClientID:     "",
	    ClientSecret: "",
	    RedirectURL:  "http://localhost:8000/oauth2callback",
	   	Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.profile",
            "https://www.googleapis.com/auth/userinfo.email",
        },
	    Endpoint: google.Endpoint,
	}


func sendOauth2Request(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	if _, ok := app.Store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	// redirect user to Google's consent page to ask for permission for the scopes specified above.
	url := conf.AuthCodeURL("uteach-login") // TODO: replace with CSRF token
	
	// redirect user to that page
    http.Redirect(w, r, url, http.StatusFound)
}

func oauth2Callback(w http.ResponseWriter, r *http.Request) {
	// handle the exchange code to initiate a transport
	authcode := r.FormValue("code")
	tok, err := conf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	response,err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	fmt.Fprint(w, string(contents), err)
}
