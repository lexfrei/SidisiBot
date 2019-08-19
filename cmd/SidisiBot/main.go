package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/lexfrei/SidisiBot/sidisilib"

	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var token string

func init() {
	viper.SetEnvPrefix("sidisi")
	err := viper.BindEnv("token")
	if err != nil {
		log.Fatalln(err)
	}

	token = viper.GetString("token")

	if token == "" {
		log.Fatalln("No token provided")
	}
}

func main() {
	var myClient *http.Client
	if sidisilib.IsThereAnyIPv6() {
		log.Println("ipv6 connection enforced")
		myClient = &http.Client{Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp6", addr)
			},
			DialTLS: func(network, addr string) (net.Conn, error) {
				return tls.Dial("tcp6", addr, &tls.Config{})
			},
		}}
	} else {
		myClient = &http.Client{}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, myClient)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil && update.InlineQuery != nil {
			go sidisilib.FuzzInline(bot, update.InlineQuery.ID, update.InlineQuery.Query)
		} else {
			go sidisilib.ResponseWithPrice(bot, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
		}

	}
}
