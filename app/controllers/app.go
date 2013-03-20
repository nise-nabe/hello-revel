package controllers

import (
	"github.com/mrjones/oauth"
	"github.com/nise-nabe/hello-revel/app/models"
	"github.com/robfig/revel"
)

var TWITTER = oauth.NewConsumer(
	"",
	"",
	oauth.ServiceProvider{
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)

type Application struct {
	*revel.Controller
}

func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) Authenticate(oauth_verifier string) revel.Result {
	user := getUser()
	if oauth_verifier != "" {
		// We got the verifier; now get the access token, store it and back to index
		accessToken, err := TWITTER.AuthorizeToken(user.RequestToken, oauth_verifier)
		if err == nil {
			user.AccessToken = accessToken
		} else {
			revel.ERROR.Println("Error connecting to twitter:", err)
		}
		return c.Redirect(Application.Index)
	}

	requestToken, url, err := TWITTER.GetRequestTokenAndUrl("http://127.0.0.1:9000/auth")
	if err == nil {
		// We received the unauthorized tokens in the OAuth object - store it before we proceed
		user.RequestToken = requestToken
		return c.Redirect(url)
	} else {
		revel.ERROR.Println("Error connecting to twitter:", err)
	}
	return c.Redirect(Application.Index)
}

func getUser() *models.User {
	return models.FindOrCreate("guest")
}
