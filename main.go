package main

import (
	"log"
	"os"
	"yc-qr-bot/pkg/agent"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var Version string

func main() {
	log.Printf("Version: %v\n", Version)
	godotenv.Load()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	// bot.ListenForWebhook(":8081")
	whInfo, _ := bot.GetWebhookInfo()
	log.Printf("whInfo: %#v\n", whInfo)
	a := agent.New(bot)
	e := echo.New()
	e.POST("/", a.HandleUpdate)
	e.Start(":" + os.Getenv("PORT"))
}
