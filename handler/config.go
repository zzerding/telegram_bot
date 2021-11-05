package handler

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BotToken string = os.Getenv("telegramBotToken")

var baseurl string = fmt.Sprintf("https://api.telegram.org/bot%s", BotToken)
var err error
var tgbot *tgbotapi.BotAPI

func init() {
	tgbot, err = tgbotapi.NewBotAPI(BotToken)
}

type message struct {
	Token string `json:"token"`
	Text  string `json:"text"`
}
