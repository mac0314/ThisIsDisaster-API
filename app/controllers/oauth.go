package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/app/routes"

	"github.com/revel/revel"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth struct {
	GorpController
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

var (
	googleConfig = &oauth2.Config{
		ClientID:     "*****-9d4ah4don21gcuv1pjrfmid3v2jc7qb5.apps.googleusercontent.com",
		ClientSecret: "*****WNnQEytMZwFKVdhO",
		RedirectURL:  "http://localhost:9000/api/oauth2/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}

	oauthStateString = "random"
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
