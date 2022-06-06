package agent

import (
	"bytes"
	"log"
	"os"
	"time"
	"yc-qr-bot/pkg/user"
	"yc-qr-bot/pkg/ydb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guregu/dynamo"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

type Agent struct {
	Bot       *tgbotapi.BotAPI
	YDBClient *dynamo.DB
}

func New(bot *tgbotapi.BotAPI) *Agent {
	return &Agent{
		Bot:       bot,
		YDBClient: ydb.New(),
	}
}

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

	}
	return c.JSON(200, nil)
}

func (a *Agent) SaveUniqueUserToYDB(update tgbotapi.Update) error {
	err := a.YDBClient.CreateTable(os.Getenv("SERVERLESS_CONTAINER_NAME")+"/users", &user.User{}).Run()
	if err != nil {
		log.Print(err)
	}
	table := a.YDBClient.Table(os.Getenv("SERVERLESS_CONTAINER_NAME") + "/users")
	count, err := table.Scan().Filter("'UserID' = ?", update.Message.From.ID).Count()
	if err != nil {
		return err
	}

	log.Printf("count by users: %v\n", count)

	if count == 0 {
		err = table.Put(&user.User{
			UserID:   update.Message.From.ID,
			Username: update.Message.From.UserName,
			Created:  time.Now(),
		}).Run()
		if err != nil {
			return err
		}
	}

	return nil
}
