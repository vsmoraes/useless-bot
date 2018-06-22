package pkg

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type (
	CustomSlackClient interface {
		PostMessage(channel string, text string, params slack.PostMessageParameters) (string, string, error)
	}

	CustomTwitterStatusService interface {
		Update(args string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error)
	}

	Handler struct {
		SlackClient   CustomSlackClient
		TwitterStatusService CustomTwitterStatusService
		Event         *slackevents.AppMentionEvent
	}

	Handleable interface {
		Handle() error
	}

	Mention struct {
		*Handler
	}
)

func NewHandler(slackCli CustomSlackClient, twitterCli CustomTwitterStatusService, ev slackevents.EventsAPIInnerEvent) (Handleable, error) {
	switch data := ev.Data.(type) {
	case *slackevents.AppMentionEvent:
		h := &Handler{SlackClient: slackCli, TwitterStatusService: twitterCli, Event: data}
		return &Mention{h}, nil
	}

	return nil, fmt.Errorf("handler not found")
}

func (h *Mention) Handle() error {
	cmd := strings.Split(strings.TrimSpace(h.Event.Text), " ")[1:][0]
	args := strings.Join(strings.Split(strings.TrimSpace(h.Event.Text), " ")[2:], " ")

	if cmd == "hey" {
		h.SlackClient.PostMessage(h.Event.Channel, "<@"+h.Event.User+">: hey there!", slack.PostMessageParameters{})
		return nil
	}

	if cmd == "tweet" {
		_, _, err := h.TwitterStatusService.Update(args, nil)
		if err != nil {
			return err
		}

		h.SlackClient.PostMessage(h.Event.Channel, "<@"+h.Event.User+">: tweet sent!\n\n"+args, slack.PostMessageParameters{})
		return nil
	}

	err := fmt.Sprintf("command not recognized: `%s`", cmd)

	h.SlackClient.PostMessage(h.Event.Channel, "<@"+h.Event.User+">: "+err, slack.PostMessageParameters{Markdown: true})
	return fmt.Errorf(err)
}
