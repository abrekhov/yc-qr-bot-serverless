package agent

import (
	"log"
	"os"
	"time"
	"yc-qr-bot/pkg/user"
	"yc-qr-bot/pkg/ydb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guregu/dynamo"
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
