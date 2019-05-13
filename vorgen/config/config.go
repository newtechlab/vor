// Package config implements simple parsing of JSON based configs for vorgen
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/juju/errgo"
)

var languages = []string{"en-US", "nb-NO"}

// A Config describes the configutation state of a answering robot.
type Config struct {
	Lang                 string   `json:"lang"`
	StartMessage         string   `json:"start_message"`
	StartMessageReply    string   `json:"start_message_reply"`
	StartMessageBadReply string   `json:"start_message_bad_reply"`
	ThanksMessage        string   `json:"thanks_message"`
	DesiredTime          int      `json:"desired_time"`
	Webhook              string   `json:"webhook"`
	SilenceTimeout       int      `json:"silence_timeout"`
	NumberVariations     int      `json:"number_variations"`
	Threads              []Thread `json:"threads"`
}

// A Thread is a series of questions that will always be asked in sequence.
type Thread []string

// Default returns an example config that may be used as a reference.
func Default() Config {
	return Config{
		Lang:                 "en-US",
		StartMessage:         "Thank you for helping us. If you agree to us storing and using this recording for training of voice models pleas say Yes.",
		StartMessageReply:    "yes",
		StartMessageBadReply: "Since you did not agree to the terms there is nothing you can help us with, thanks anyway.",
		ThanksMessage:        "Thanks for calling, please remember to make 3 calls from different environments but the same phone",
		DesiredTime:          20,
		Webhook:              "https://example.com/api/callback",
		SilenceTimeout:       2,
		NumberVariations:     10,
		Threads: []Thread{
			Thread{"Do you like cats or dogs the most?", "Why is that?"},
			Thread{"Why do you help us?", "Would you ever do it again?"},
			Thread{"How is the weather like today?", "How does that compare to yesterday"},
			Thread{"What is yout favourite colors?"},
			Thread{"Please count from one to twenty."},
			Thread{"Do you like smokers?"},
			Thread{"Do you like to swim?"},
			Thread{"Describe your best vaccation."},
			Thread{"Are you at work?"},
		},
	}
}

func (c Config) String() string {
	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "{error encoding config to JSON}"
	}
	return string(buf)
}

// LoadConfig tries to parse a JSON config from the reader and validates it. Failure of
// either step will return nil and an error.
func LoadConfig(r io.Reader) (c Config, err error) {
	c = Default()
	dec := json.NewDecoder(r)
	err = dec.Decode(&c)
	return c, errgo.Mask(err)
}

func (c Config) validate() error {
	if err := c.checkLang(); err != nil {
		return errgo.Mask(err)
	}
	if c.StartMessage == "" {
		return errgo.New("'start_message' must be set to the initial message to be played to the user")
	}
	if c.StartMessageReply == "" {
		return errgo.New("'start_message_reply' must be specified as the expression the user must say to indicate agreement to the terms/start_message. Typically YES in the given language.")
	}
	if c.Threads == nil || len(c.Threads) == 0 {
		return errgo.New("'threads' must specified as an array of array of strings")
	}
	for i, t := range c.Threads {
		if t == nil || len(t) == 0 {
			return errgo.New("error in thread no " + strconv.Itoa(i) + ": each thread must be an array of strings with at least one question to ask.")
		}
	}
	return nil
}

func (c Config) checkLang() error {
	for _, l := range languages {
		if l == c.Lang {
			return nil
		}
	}
	return errgo.New("unknown language code: " + c.Lang + ", accepted options are: " + fmt.Sprint(languages))
}
