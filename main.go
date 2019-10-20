package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {
	cfg := config{
		oauthID:     os.Getenv("oauthID"),
		oauthSecret: os.Getenv("oauthSecret"),
	}
	lambda.Start(func() error {
		return handler(cfg)
	})
}

type config struct {
	oauthID     string
	oauthSecret string
}

func handler(cfg config) error {
	s := newSession(cfg.oauthID, cfg.oauthSecret)
	s.login()
	s.poll()
	s.download()
	return s.err
}
