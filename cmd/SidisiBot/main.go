package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lexfrei/goscgp/parser"
)

var token string

func init() {
	viper.SetEnvPrefix("sidisi")
	viper.BindEnv("token")

	token = viper.GetString("token")

	if token == "" {
		log.Fatalln("No token provided")
	}
}

func main() {
	// // Enforce ipv6 connection
	// myClient := &http.Client{Transport: &http.Transport{
	// 	Dial: func(network, addr string) (net.Conn, error) {
	// 		return net.Dial("tcp6", addr)
	// 	},
	// 	DialTLS: func(network, addr string) (net.Conn, error) {
	// 		return tls.Dial("tcp6", addr, &tls.Config{})
	// 	},
	// }}

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
			var articles []interface{}

			card := fuzzCard(update.InlineQuery.Query)

			msg := tgbotapi.NewInlineQueryResultArticle(
				update.InlineQuery.ID,
				card,
				card)
			articles = append(articles, msg)

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    false,
				CacheTime:     0,
				Results:       articles,
			}
			_, err := bot.AnswerInlineQuery(inlineConfig)
			if err != nil {
				log.Println(err)
			}

		} else {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, scgPrices(update.Message.Text))
			msg.ParseMode = "markdown"
			fmt.Println(update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}

}

func fuzzCard(card string) string {
	log.Println("Trying to find:", card)
	var c Card
	u, err := url.Parse("https://api.scryfall.com/cards/named")
	if err != nil {
		log.Println(err)
	}

	q := u.Query()
	q.Set("fuzzy", card)

	u.RawQuery = q.Encode()

	log.Println(u.String())

	res, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
	}

	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&c)

	log.Println(c.Name)

	return c.Name
}

func scgPrices(card string) string {
	siteURL, err := url.Parse("http://www.starcitygames.com/results?&switch_display=1")
	if err != nil {
		log.Println(err)
	}

	q := siteURL.Query()
	q.Set("name", card)
	siteURL.RawQuery = q.Encode()

	log.Println(siteURL.String())

	c := &http.Client{}

	result, err := parser.DoRequest(*siteURL, c)
	if err != nil {
		log.Println(err)
	}

	log.Println("Got " + strconv.Itoa(len(result)) + " cards")

	var str string

	for _, v := range result {
		str = str + v.String() + "\n"
	}

	return str
}

type Card struct {
	Name string `json:"name"`
}
