package agent

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"

	_ "image/jpeg"
	_ "image/png"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/liyue201/goqr"
	"github.com/skip2/go-qrcode"
)

func (a *Agent) HandleUpdate(c echo.Context) error {
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		panic(err)
	}

	if update.Message != nil {
		a.SaveUniqueUserToYDB(update)
	} else {
		return c.JSON(204, nil)
	}

	if update.Message.Text != "" {
		err := a.EncodeContent(update)
		log.Print(err)
	}
	if update.Message.Photo != nil || update.Message.Document != nil {
		err := a.DecodeContent(update)
		log.Print(err)
	}
	return c.JSON(200, nil)
}

func (a *Agent) EncodeContent(update tgbotapi.Update) error {
	pngBytes, err := qrcode.Encode(update.Message.Text, qrcode.Low, 256)
	if err != nil {
		log.Print(err)
		return err
	}
	msg := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FileReader{
		Name:   "QR",
		Reader: bytes.NewBuffer(pngBytes),
	})
	msg.ReplyToMessageID = update.Message.MessageID
	a.Bot.Send(msg)
	return nil

}

func (a *Agent) DecodeContent(update tgbotapi.Update) error {
	lenPhoto := len(update.Message.Photo)
	downloadLink, err := a.Bot.GetFileDirectURL(update.Message.Photo[lenPhoto-1].FileID)
	if err != nil {
		log.Print(err)
		return err
	}
	fileContent, err := http.Get(downloadLink)
	if err != nil {
		log.Print(err)
		return err
	}
	defer fileContent.Body.Close()

	img, _, err := image.Decode(fileContent.Body)
	if err != nil {
		log.Print(err)
		return err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Print(err)
		return err
	}
	var finalMessage string
	for _, qrCode := range qrCodes {
		finalMessage += fmt.Sprintf("%s\n", qrCode.Payload)
	}
	msg := tgbotapi.NewMessage(update.Message.From.ID, finalMessage)
	msg.ReplyToMessageID = update.Message.MessageID
	a.Bot.Send(msg)
	return nil
}
