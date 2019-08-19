package sidisilib

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lexfrei/SidisiBot/types"
	"github.com/lexfrei/goscgp/parser"
)

func GetSCGPrices(card string) string {
	siteURL, err := url.Parse("http://www.starcitygames.com/results?&switch_display=1")
	if err != nil {
		log.Println("can't parse url:", err)
	}

	q := siteURL.Query()
	q.Set("name", card)
	siteURL.RawQuery = q.Encode()

	c := &http.Client{}

	result, err := parser.DoRequest(*siteURL, c)
	if err != nil {
		log.Println("can't do request: ", err)
	}

	if len(result) == 0 {
		return "*Zero cards found*"
	}

	var str string
	var cardsParsed int
	for _, v := range result {
		if len(str+v.String()+"\n") > 4000 {
			return fmt.Sprintf("*Too many cards. %d/%d are shown*\n\n%s", cardsParsed, len(result), str)
		}
		cardsParsed++
		str = str + v.String() + "\n"
	}

	return str
}

// ResponseWithPrice comment placeholder
func ResponseWithPrice(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string) {

	// edit := tgbotapi.NewEditMessageText(chatID, messageID, "kek!")
	// _, err := bot.Send(edit)
	// if err != nil {
	// 	log.Println(err)
	// }

	msg := tgbotapi.NewMessage(chatID, GetSCGPrices(text))
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = messageID
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("can't do response:", err)
	}
}

// FuzzInline comment placeholder
func FuzzInline(bot *tgbotapi.BotAPI, qID string, text string) {
	if len(text) == 0 {
		return
	}
	var articles []interface{}

	card := fuzzCard(text)
	if len(card) == 0 {
		return
	}

	msg := tgbotapi.NewInlineQueryResultArticle(
		qID,
		card,
		card)
	articles = append(articles, msg)

	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: qID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       articles,
	}
	_, err := bot.AnswerInlineQuery(inlineConfig)

	if err != nil {
		log.Println("can't do inline response", err)
	}
}

func fuzzCard(card string) string {
	var c types.Card
	u, err := url.Parse("https://api.scryfall.com/cards/named")
	if err != nil {
		log.Println(err)
	}

	q := u.Query()
	q.Set("fuzzy", card)

	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		log.Println("can't do get", err)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&c)

	if err != nil {
		log.Println("can't decode", err)
		return ""
	}

	return c.Name
}

func IsThereAnyIPv6() bool {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.To4() != nil {
				return true
			}
		}
	}
	return false
}
