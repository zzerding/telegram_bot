package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"telegramBot/handler/internal/bingwallpaper"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var tokenError = errors.New("token error")

func BingWall(c *gin.Context) {
	data, err := bingwallpaper.GetCache(bingwallpaper.UrlBingServer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, data)

}
func Ip(c *gin.Context) {
	c.String(http.StatusOK, c.ClientIP())
}
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
func SendMsg(c *gin.Context) {
	m := new(message)
	err := c.ShouldBind(m)
	if err != nil {
		badRequest(c, err)
		return
	}
	req, _ := json.Marshal(m)
	log.Printf("%s", req)
	id, err := idDecode(m.Token)
	if err != nil {
		log.Println(err)
		badRequest(c, tokenError)
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(id, err)
		badRequest(c, tokenError)
		return
	}
	msg := tgbotapi.NewMessage(idInt, m.Text)
	_, err = tgbot.Send(msg)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func helpMsg(id int) tgbotapi.MessageConfig {
	text := `
这个 bot 使用的是vercel服务，主要用来推送消息的，因为国内有的服务器访问不了 telegram
使用方法：
	1. 取得自己的 token /mytoken
	2. curl -d "token=yourtoken&text=your msg" https://zzerd.vercel.app/sendmsg  
	可以使用你喜欢的任意语言 post 这个网址支持 json 与 form 格式
	`
	return tgbotapi.NewMessage(int64(id), text)
}
func tokenMsg(id int) (msg tgbotapi.MessageConfig, err error) {
	token, err := idEncode(fmt.Sprint(id))
	if err != nil {
		return
	}
	msg = tgbotapi.NewMessage(int64(id), token)
	return
}
func errorMsg(id int) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(int64(id), "server 500")
}
func UseHook(c *gin.Context) {
	callBack := new(tgbotapi.CallbackQuery)
	err := c.ShouldBindJSON(callBack)
	if err != nil {
		badRequest(c, err)
		return
	}
	msg := callBack.Message
	respose, err := json.Marshal(callBack)
	log.Println(string(respose))
	if msg == nil {
		badRequest(c, err)
		return
	}
	// chatID := msg.Chat.ID
	userID := msg.From.ID
	command := msg.Command()
	log.Println("input command", command)
	switch command {
	case "start":
		tgbot.Send(helpMsg(userID))
		log.Println("command start")
	case "help":
		tgbot.Send(helpMsg(userID))
		log.Println("command help")
	case "mytoken":
		sendMsg, err := tokenMsg(userID)
		if err != nil {
			tgbot.Send(errorMsg(userID))
			return
		}
		tgbot.Send(sendMsg)
		log.Println("command mytoken")
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
