package handlers

import (
	"fmt"
	"github.com/barakb/github-branch/session"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
	"os"
)

var oauthConf = &oauth2.Config{
	ClientID:     os.Getenv("ClientID"),
	ClientSecret: os.Getenv("ClientSecret"),
	// Comma separated list of scopes
	// select level of access you want https://developer.github.com/v3/oauth/#scopes
	Scopes:   []string{"user:email, repo"},
	Endpoint: githuboauth.Endpoint,
}

var oauthStateString = "arandomstring"

type GithubAuthHandler struct {
	next http.Handler
}

func (h GithubAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GithubAuthHandler\n")
	sess := session.GlobalSessions.SessionStart(w, r)
	client := sess.Get("*github.client")
	if client == nil {
		fmt.Printf("sess is %#v\n", sess)
		fmt.Printf("r.RequestURI is %#v\n", r.RequestURI)
		sess.Set("redirect", r.RequestURI)
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(h http.Handler) http.Handler {
	return GithubAuthHandler{next: h}
}
