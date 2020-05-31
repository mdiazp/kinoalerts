package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kimrgrey/go-telegram"
	"github.com/rs/xid"
)

type environment struct {
	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN"`

	IntervalTimeToQuery      int `envconfig:"INTERVAL_TIME_TO_QUERY"`
	IntervalTimeToNotificate int `envconfig:"INTERVAL_TIME_TO_NOTIFICATE"`
}

// Alert ...
type Alert struct {
	ID            string
	Name          string   `json:"name"`
	RegexPatterns []string `json:"regexPatterns"`
}

// Resource ...
type Resource struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	ChatID string  `json:"chatID"`
	Alerts []Alert `json:"alerts"`
}

// Config ...
type Config struct {
	Resources []Resource `json:"resources"`
}

var (
	env            environment
	telegramClient *telegram.Client
	notifications  map[string]time.Time = map[string]time.Time{}
	config         Config
)

func initTelegramClient() {
	telegramClient = telegram.NewClient(env.TelegramBotToken)
}

func printTelegramAccount() {
	fmt.Printf("%v\n", telegramClient.GetMe())
}

func sendMessage(chatID string, message string) {
	params := url.Values{
		"chat_id": []string{chatID},
		"text":    []string{message},
	}

	var v map[string]interface{}
	telegramClient.Call("sendMessage", params, &v)

	if ok, good := v["ok"].(bool); !ok || !good {
		fmt.Printf("call: %v\n", v)
	}
}

func checkPattern(s string, pattern string) bool {
	// ok := strings.Contains(s, pattern)
	ok, _ := regexp.MatchString(pattern, s)
	return ok
}

func lookOnResource(resource Resource) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r, e := http.Post(
		resource.URL,
		"text/plain; charset=utf-8",
		bytes.NewReader([]byte{}),
	)

	if e != nil {
		return fmt.Errorf("Unable to send request: %v", e)
	}

	rb, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return fmt.Errorf("Unable to unmarshal response: %v", e)
	}

	s := string(rb)
	s = strings.ToLower(s)

	notification := ""

	for _, alert := range resource.Alerts {
		founded := false
		for _, pattern := range alert.RegexPatterns {
			if checkPattern(s, pattern) {
				founded = true
				break
			}
		}

		if founded {
			notificate := true
			if last, ok := notifications[alert.ID]; ok {
				notificate = last.Add(time.Second * time.Duration(env.IntervalTimeToNotificate)).Before(time.Now())
			}

			if notificate {
				newNot := fmt.Sprintf("Hay %s!!!\nLink: %s", alert.Name, resource.URL)

				notification = fmt.Sprintf("%s\n%s", notification, newNot)

				notifications[alert.ID] = time.Now()
			}
		}
	}

	if notification != "" {
		sendMessage(resource.ChatID, notification)
	}

	return nil
}

func start(resource Resource) {
	for i := range resource.Alerts {
		resource.Alerts[i].ID = xid.New().String()
	}

	for {
		e := lookOnResource(resource)
		if e != nil {
			fmt.Printf("Error looking on resource %s:\n\t%v\n", resource.Name, e)
		}

		time.Sleep(time.Second * time.Duration(env.IntervalTimeToQuery))
	}
}

func readResources() error {
	f, e := os.Open("resources.json")
	if e != nil {
		return fmt.Errorf("Unable to read resource.json file: %v", e)
	}

	bconfig, e := ioutil.ReadAll(f)
	if e != nil {
		return fmt.Errorf("Unable to read resource.json file: %v", e)
	}

	e = json.Unmarshal(bconfig, &config)
	if e != nil {
		return fmt.Errorf("Unable to read resource.json file: %v", e)
	}

	return nil
}

func main() {
	we := func(e error) {
		if e == nil {
			return
		}
		fmt.Printf("ERROR: %v", e)
		os.Exit(1)
	}
	var e error

	e = envconfig.Process("", &env)
	we(e)

	e = readResources()
	we(e)

	// fmt.Printf("environment: %v\n", env)

	initTelegramClient()

	stop := make(chan bool)

	for _, resource := range config.Resources {
		go start(resource)
	}

	// Capture os signal to stop
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch

		stop <- true
	}()

	<-stop
}
