package pkg

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	MockSlackClient struct {
		mock.Mock
		TimesCalled int
		Channels map[int]string
		Texts map[int]string
		Params map[int]slack.PostMessageParameters
	}
	MockStatusService struct {
		mock.Mock
		TimesCalled int
		Args map[int]string
		Params map[int]*twitter.StatusUpdateParams
	}
)

func (m *MockSlackClient) PostMessage(channel string, text string, params slack.PostMessageParameters) (string, string, error) {
	index := m.TimesCalled - 1
	if index < 0 {
		index = 0
		m.TimesCalled = 0
	}

	if m.Channels == nil {
		m.Channels = make(map[int]string)
	}

	if m.Texts == nil {
		m.Texts = make(map[int]string)
	}

	if m.Params == nil {
		m.Params = make(map[int]slack.PostMessageParameters)
	}

	m.TimesCalled = m.TimesCalled + 1
	m.Channels[index] = channel
	m.Texts[index] = text
	m.Params[index] = params

	return "", "", nil
}

func (m *MockStatusService) Update(args string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error) {
	index := m.TimesCalled - 1
	if index < 0 {
		index = 0
		m.TimesCalled = 0
	}

	if m.Args == nil {
		m.Args= make(map[int]string)
	}

	if m.Params == nil {
		m.Params = make(map[int]*twitter.StatusUpdateParams)
	}

	m.TimesCalled = m.TimesCalled + 1
	m.Args[index] = args
	m.Params[index] = params

	return nil, nil, nil
}

func TestShouldCreateNewMentionHandler(t *testing.T) {
	slackCli := slack.Client{}
	twitterCli := MockStatusService{}
	event := slackevents.EventsAPIInnerEvent{
		Data: &slackevents.AppMentionEvent{},
	}

	handler, err := NewHandler(&slackCli, &twitterCli, event)

	assert.IsType(t, &Mention{}, handler)
	assert.Nil(t, err)
}

func TestShouldReturnErrorWhenEventIsNotRecognized(t *testing.T) {
	slackCli := MockSlackClient{}
	twitterCli := MockStatusService{}
	event := slackevents.EventsAPIInnerEvent{}

	handler, err := NewHandler(&slackCli, &twitterCli, event)

	assert.Nil(t, handler)
	assert.Error(t, err)
	assert.Equal(t, 0, slackCli.TimesCalled)
	assert.Equal(t, 0, twitterCli.TimesCalled)
}

func TestShouldHandleMentionNotRecognized(t *testing.T) {
	channel := "fake-channel"
	user := "FAKE-ID"

	slackCli := MockSlackClient{}
	twitterCli := MockStatusService{}
	event := slackevents.EventsAPIInnerEvent{
		Data: &slackevents.AppMentionEvent{
			Channel: channel,
			User: user,
			Text: "<@FAKE-BOT-ID>: not-recognized-command foo bar",
		},
	}

	handler, _ := NewHandler(&slackCli, &twitterCli, event)
	err := handler.Handle()

	expectedText := fmt.Sprintf("<@%s>: %s", user, err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, 1, slackCli.TimesCalled)
	assert.Equal(t, slackCli.Channels[0], channel)
	assert.Equal(t, slackCli.Texts[0], expectedText)
	assert.Equal(t, slackCli.Params[0], slack.PostMessageParameters{Markdown:true})
	assert.Equal(t, 0, twitterCli.TimesCalled)
}

func TestShouldHandleMentionHey(t *testing.T) {
	channel := "fake-channel"
	user := "FAKE-ID"
	text := "<@"+user+">: hey there!"

	slackCli := MockSlackClient{}
	twitterCli := MockStatusService{}
	event := slackevents.EventsAPIInnerEvent{
		Data: &slackevents.AppMentionEvent{
			Channel: channel,
			User: user,
			Text: "<@FAKE-BOT-ID>: hey",
		},
	}

	handler, _ := NewHandler(&slackCli, &twitterCli, event)
	err := handler.Handle()

	assert.Nil(t, err)
	assert.Equal(t, 1, slackCli.TimesCalled)
	assert.Equal(t, slackCli.Channels[0], channel)
	assert.Equal(t, slackCli.Texts[0], text)
	assert.Equal(t, slackCli.Params[0], slack.PostMessageParameters{})
	assert.Equal(t, 0, twitterCli.TimesCalled)
}

func TestShouldHandleMentionTweet(t *testing.T) {
	channel := "fake-channel"
	user := "FAKE-ID"
	tweet := "this is a fake tweet"
	text := "<@"+user+">: tweet sent!\n\n" + tweet

	slackCli := MockSlackClient{}
	twitterCli := MockStatusService{}

	event := slackevents.EventsAPIInnerEvent{
		Data: &slackevents.AppMentionEvent{
			Channel: channel,
			User: user,
			Text: "<@FAKE-BOT-ID>: tweet " + tweet,
		},
	}

	handler, _ := NewHandler(&slackCli, &twitterCli, event)
	err := handler.Handle()
	assert.Nil(t, err)
	assert.Equal(t, 1, slackCli.TimesCalled)
	assert.Equal(t, slackCli.Channels[0], channel)
	assert.Equal(t, slackCli.Texts[0], text)
	assert.Equal(t, slackCli.Params[0], slack.PostMessageParameters{})
}