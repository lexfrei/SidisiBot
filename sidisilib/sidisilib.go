package sidisilib

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	scryfall "github.com/BlueMonday/go-scryfall"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lexfrei/goscgp/parser"
)

var reEngLan = regexp.MustCompile(`^[A-Za-z\s]+$`)

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

	ctx := context.Background()

	client, err := scryfall.NewClient()
	if err != nil {
		log.Println("can't create scryfall client:", err)
		return
	}

	result, err := client.SearchCards(ctx, text, scryfall.SearchCardsOptions{Unique: scryfall.UniqueModeCards, IncludeMultilingual: !reEngLan.Match([]byte(text))})
	if err != nil {
		log.Printf("can't fuzz search for \"%s\": %s", text, err)
		return
	}

	if len(result.Cards) < 1 {
		return
	}

	msg := tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    qID,
		Title: result.Cards[0].Name,
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: result.Cards[0].Name,
		},
	}
	if result.Cards[0].ImageURIs != nil {
		msg.ThumbURL = result.Cards[0].ImageURIs.ArtCrop
		msg.ThumbHeight = 50
	}

	articles = append(articles, msg)

	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: qID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       articles,
	}
	_, err = bot.AnswerInlineQuery(inlineConfig)
	if err != nil {
		log.Println("can't do inline response", err)
	}
}
