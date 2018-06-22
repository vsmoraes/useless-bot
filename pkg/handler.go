package pkg

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
)

type (
	IndexHandler struct {
		BotToken              string
		VerificationToken     string
		TwitterConsumerKey    string
		TwitterConsumerSecret string
		TwitterAccessKey      string
		TwitterAccessSecret   string
	}
)

func (h *IndexHandler) ProcessEvent(c echo.Context) error {
	slackClient := slack.New(h.BotToken)
	twtConfig := oauth1.NewConfig(h.TwitterConsumerKey, h.TwitterConsumerSecret)
	twtToken := oauth1.NewToken(h.TwitterAccessKey, h.TwitterAccessSecret)
	twitterClient := twitter.NewClient(twtConfig.Client(oauth1.NoContext, twtToken))

	ev := &Events{
		SlackClient:       slackClient,
		TwitterClient:     twitterClient,
		VerificationToken: h.VerificationToken,
	}

	return ev.Handle(c)
}
