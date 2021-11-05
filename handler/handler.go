package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}

func SendMsg(c *gin.Context) {
	m := new(message)
	err := c.ShouldBind(m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id := idDecode(m.Token)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	msg := tgbotapi.NewMessage(idInt, m.Text)
	_, err = tgbot.Send(msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func helpMsg(id int) tgbotapi.MessageConfig {
	text := `help`
	return tgbotapi.NewMessage(int64(id), text)
}
func tokenMsg(id int) tgbotapi.MessageConfig {
	token := idEncode(fmt.Sprint(id))
	return tgbotapi.NewMessage(int64(id), token)
}
func errorMsg(id int) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(int64(id), "server 500")
}
func UseHook(c *gin.Context) {
	callBack := new(tgbotapi.CallbackQuery)
	err := c.ShouldBindJSON(callBack)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request error"})
		return
	}
	log.Println(callBack)
	msg := callBack.Message
	if msg == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message not found"})
		return
	}
	// chatID := msg.Chat.ID
	userID := msg.From.ID
	command := msg.Command()
	switch command {
	case "/star":
		tgbot.Send(helpMsg(userID))
	case "/help":
		tgbot.Send(helpMsg(userID))
	case "/token":
		tgbot.Send(tokenMsg(userID))
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
func SetHook(c *gin.Context) {
	webHook := tgbotapi.NewWebhook(fmt.Sprintf("%s/%s", "https://zzerd.vercel.app/api", BotToken))
	_, err := tgbot.SetWebhook(webHook)
	if err != nil {
		log.Panicln(err)
	}
}
