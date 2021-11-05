package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"telegramBot/handler"
)

var (
	app *gin.Engine
)

func registerRouter(r *gin.RouterGroup) {
	r.GET("/setHook", handler.SetHook)
	r.GET("/ping", handler.Ping)
	r.POST("/sendmsg", handler.SendMsg)
	r.POST(handler.BotToken, handler.UseHook)
}

// init gin app
func init() {
	app = gin.New()

	// Handling routing errors
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no route"})
	})
	// must /api/xxx
	r := app.Group("/api")

	// register route
	registerRouter(r)
}

// entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
