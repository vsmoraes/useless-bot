package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type (
	Events struct {
		SlackClient       *slack.Client
		TwitterClient     *twitter.Client
		VerificationToken string
	}

	ErrorReponse struct {
		Message string `json:"message"`
	}

	ChallengeResponse struct {
		Challenge string `json:"challenge"`
	}
)

func (e *Events) Handle(c echo.Context) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	ev := buf.String()

	event, err := slackevents.ParseEvent(json.RawMessage(ev), slackevents.OptionVerifyToken(&slackevents.TokenComparator{e.VerificationToken}))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorReponse{Message: err.Error()})
	}

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(ev), &r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrorReponse{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, &ChallengeResponse{Challenge: r.Challenge})
	}

	if event.Type == slackevents.CallbackEvent {
		handler, err := NewHandler(e.SlackClient, e.TwitterClient.Statuses, event.InnerEvent)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrorReponse{Message: err.Error()})
		}

		handler.Handle()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrorReponse{Message: err.Error()})
		}
	}

	return c.NoContent(http.StatusOK)
}
