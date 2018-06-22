package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/vsmoraes/useless-bot/pkg"
)

type (
	EnvConfig struct {
		Port                  string `envconfig:"PORT" default:"3000"`
		BotToken              string `envconfig:"BOT_TOKEN" required:"true"`
		VerificationToken     string `envconfig:"VERIFICATION_TOKEN" required:"true"`
		TwitterConsumerKey    string `envconfig:"TWITTER_CONSUMER_KEY" required:"true"`
		TwitterConsumerSecret string `envconfig:"TWITTER_CONSUMER_SECRET" required:"true"`
		TwitterAccessKey      string `envconfig:"TWITTER_ACCESS_KEY" required:"true"`
		TwitterAccessSecret   string `envconfig:"TWITTER_ACCESS_SECRET" required:"true"`
	}
)

func main() {
	var env EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	indexHandler := &pkg.IndexHandler{
		BotToken:              env.BotToken,
		VerificationToken:     env.VerificationToken,
		TwitterConsumerKey:    env.TwitterConsumerKey,
		TwitterConsumerSecret: env.TwitterConsumerSecret,
		TwitterAccessKey:      env.TwitterAccessKey,
		TwitterAccessSecret:   env.TwitterAccessSecret,
	}

	e := echo.New()
	e.POST("/", indexHandler.ProcessEvent)

	e.Logger.Fatal(e.Start(":3000"))
}
