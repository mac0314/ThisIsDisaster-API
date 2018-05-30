package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"server/app/routes"

	"github.com/revel/revel"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type OAuth struct {
	GorpController
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

var (
	googleConfig = &oauth2.Config{
		ClientID:     "791663390275-9d4ah4don21gcuv1pjrfmid3v2jc7qb5.apps.googleusercontent.com",
		ClientSecret: "2sOMmEuSWNnQEytMZwFKVdhO",
		RedirectURL:  "http://localhost:9000/api/oauth2/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}

	oauthStateString = "random"

	facebookConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{},
		Endpoint:     facebook.Endpoint,
		RedirectURL:  "http://localhost:9000/api/oauth2/facebook/callback",
	}

	githubConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9000/api/oauth2/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func (c OAuth) GoogleAuth(code string) revel.Result {

	if len(code) > 0 {
		tok, err := googleConfig.Exchange(oauth2.NoContext, code)
		log.Println(tok)
		if err != nil {
			log.Println("err is", err)
			return c.Redirect(routes.App.Index())
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
		defer response.Body.Close()

		fmt.Println(response)
		fmt.Println("2")
		return c.Redirect(routes.App.Index())
	} else {
		url := googleConfig.AuthCodeURL("")
		log.Println(url)
		fmt.Println("3")
		return c.Redirect(url)

	}
}

func (c OAuth) Callback() revel.Result {

	state := c.Params.Form.Get("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.Redirect(routes.App.Index())
	}

	code := c.Params.Form.Get("code")
	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return c.Redirect(routes.App.Index())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Printf("Content: %s\n", contents)
	return c.Redirect(routes.App.Index())
}

func (c OAuth) FacebookAuth(code string) revel.Result {
	if len(code) > 0 {
		token, err := facebookConfig.Exchange(oauth2.NoContext, code)

		if err != nil {
			log.Println("err is", err)

			return c.Redirect(routes.App.Index())
		}

		response, err := http.Get("https://graph.facebook.com/v2.4/me?field=id,picture,email&access_token=" + token.AccessToken)

		defer response.Body.Close()
		if err != nil {
			log.Println(err)
		}

	} else {
		url := facebookConfig.AuthCodeURL("")

		return c.Redirect(url)
	}
	return c.Redirect(routes.App.Index())
}

func (c OAuth) GithubAuth(code string) revel.Result {
	if len(code) > 0 {
		token, err := githubConfig.Exchange(oauth2.NoContext, code)

		if err != nil {
			log.Println("err is", err)

			return c.Redirect(routes.App.Index())
		}

		log.Println("https://api.github.com/?access_token=" + token.AccessToken)

		response, err := http.Get("https://api.github.com/user?access_token=" + token.AccessToken)

		defer response.Body.Close()
		if err != nil {
			log.Println(err)
		}

		return c.Redirect(routes.App.Index())
	} else {
		url := githubConfig.AuthCodeURL("")

		return c.Redirect(url)
	}

	return c.Redirect(routes.App.Index())
}

func (c OAuth) randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (c OAuth) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect(routes.App.Index())
}
