package agent

import (
	"bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

type Agent struct {
	Bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *Agent {
	return &Agent{
		Bot: bot,
	}
}

func (a *Agent) HandleUpdate(c echo.Context) error {
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		panic(err)
	}
	pngBytes, err := qrcode.Encode(update.Message.Text, qrcode.Low, 256)
	if err != nil {
		panic(err)
	}
	msg := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FileReader{
		Name:   "QR",
		Reader: bytes.NewBuffer(pngBytes),
	})
	msg.ReplyToMessageID = update.Message.MessageID
	a.Bot.Send(msg)
	return c.JSON(200, nil)
}
