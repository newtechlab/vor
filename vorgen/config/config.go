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
	// Language for both reading up texts and understanding the user
	Lang string `json:"lang"`
	// The initial Message played when the user calls, typically a question
	// to ask for acceptance to use the data in addition to any other information
	StartMessage string `json:"start_message"`
	// The reply the user have to give to the initial message to proceed with the
	// recording, typically somehting like "YES"
	StartMessageReply string `json:"start_message_reply"`
	// What the robot should say if the user did not say StartMessageReply
	StartMessageBadReply string `json:"start_message_bad_reply"`
	// Message that should be played at the end of the questions, when the user
	// has replied as desired.
	ThanksMessage string `json:"thanks_message"`
	// How long a total recording is desired, the robot will keep asking questions
	// (provided enough are defined) until a recording of this length has been achieved
	DesiredTime int `json:"desired_time"`
	// URL to the vorgserve (r) that should be sent the recordings.
	Webhook string `json:"webhook"`
	// How long a silence interval to wait before asking the next question
	SilenceTimeout int `json:"silence_timeout"`
	// How many different series of questions should be generated, a reasonable value
	// may be 3 or 5. Performance of Twillio scales badly with large numbers,
	NumberVariations int `json:"number_variations"`
	// Series of questions that should be asked by the robot (see example below). Note
	// That the order of thre threads will be randomized (based on NumberVariations),
	// but that the order of questions inside each thread will be the same. This allows
	// for threads of questions that will help the user to think about the stuff that
	// will be asked next.
	Threads []Thread `json:"threads"`
}

// A Thread is a series of questions that will always be asked in sequence.
type Thread []string

// // Default returns an example config that may be used as a reference.
// func Default() Config {
// 	return Config{
// 		Lang:                 "en-US",
// 		StartMessage:         "Thank you for helping us. If you agree to us storing and using this recording for training of voice models pleas say Yes.",
// 		StartMessageReply:    "yes",
// 		StartMessageBadReply: "Since you did not agree to the terms there is nothing you can help us with, thanks anyway.",
// 		ThanksMessage:        "Thanks for calling, please remember to make 3 calls from different environments but the same phone",
// 		DesiredTime:          180,
// 		Webhook:              "https://example.com/api/callback",
// 		SilenceTimeout:       2,
// 		NumberVariations:     10,
// 		Threads: []Thread{
// 			Thread{"A thread is a series of questions.", "That will always be asked in sequence.", "The order between different threads will be randomized though."},
// 			Thread{"Do you like cats or dogs the most?", "Why is that?"},
// 			Thread{"Why do you help us?", "Would you ever do it again?"},
// 			Thread{"How is the weather like today?", "How does that compare to yesterday"},
// 			Thread{"What is yout favourite colors?"},
// 			Thread{"Please count from one to twenty."},
// 			Thread{"Do you like smokers?"},
// 			Thread{"Do you like to swim?"},
// 			Thread{"Describe your best vaccation."},
// 			Thread{"Are you at work?"},
// 		},
// 	}
// }
// NORWEGIAN Default returns an example config that may be used as a reference.
func Default() Config {
	return Config{
		Lang:                 "nb-NO",
		StartMessage:         "Takk for at du vil bidra. Samtykker du til at vi bruker opptaket fra din samtale til å trene og validere en modell for stemmeidentifikasjon?",
		StartMessageReply:    "ja",
		StartMessageBadReply: "Det er påkrevd at du samtykker for at vi skal kunne gjøre et opptak av din samtale. Takk for at du ringte.",
		ThanksMessage:        "Takk for at du ringer. For å kunne teste systemet best mulig trenger vi opptak av tre samtaler fra deg, helst fra tre ulike steder.",
		DesiredTime:          30,
		Webhook:              "https://vorno.newtechlab.wtf",
		SilenceTimeout:       4,
		NumberVariations:     10,
		Threads: []Thread{
			Thread{"Hva gjorde du i sommerferien?", "Liker du best sommer eller vinter, og hvorfor det?"},
			Thread{"Hvorfor valgte du å bidra ved å ringe inn her?", "Bidrar du med andre frivillige aktiviteter eller organisasjoner?"},
			Thread{"Beskriv hvordan været er der du er i dag?", "Hva synes du om dagen vær?"},
			Thread{"Hva liker du best å gjøre i fritiden?"},
			Thread{"Tell fra null til tjue."},
			Thread{"Hva synes du om åpne kontorlandskap?"},
			Thread{"Si navnet på alle månedene i året"},
			Thread{"Beskriv din drømmeferie."},
			Thread{"Hva er det beste med jobben din?", "Hva liker du best å spise til lunsj og hvorfor?"},
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
