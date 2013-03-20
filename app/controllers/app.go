package controllers

import (
	"encoding/json"
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
	user := getUser()
	if user.AccessToken != nil {
		resp, err := TWITTER.Get(
			"https://api.twitter.com/1.1/statuses/home_timeline.json",
			map[string]string{"count": "10"},
			user.AccessToken)
		if err != nil {
			revel.ERROR.Println(err)
			return c.Render()
		}
		defer resp.Body.Close()

		// Extract the mention text.
		tweets := []struct {
			Text string `json:text`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&tweets)

		return c.Render(tweets)
	}
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
