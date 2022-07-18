package telegrambot

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
	"tool/telegrambot/model"
)

var robotInstance *TelegramRobot
var telegramRobotOnce sync.Once

type TelegramBot interface {
	Name() string
	SetChatID(chatID int64) int64
	SendMessage(respMessage model.Message) error
}

type TelegramRobot struct {
	token  string
	chatID int64
	Robot  *tgbotapi.BotAPI
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTelegramRobot(token string, chatID int64) TelegramBot {
	if robotInstance != nil {
		return robotInstance
	}

	telegramRobotOnce.Do(func() {
		robotInstance = &TelegramRobot{}
		robotInstance.token = token
		robotInstance.chatID = chatID
		robotInstance.ConnectBotServer(token)
		robotInstance.ctx, robotInstance.cancel = context.WithCancel(context.Background())
	})

	return robotInstance
}

func GetTelegramInstance() TelegramBot {
	if robotInstance != nil && robotInstance.Robot != nil {
		return robotInstance
	}
	return nil
}

func (bot *TelegramRobot) ConnectBotServer(token string) {
	if bot.Robot == nil {
		botConnect, err := tgbotapi.NewBotAPI(token)
		if err != nil {

			return
		}
		bot.Robot = botConnect
	}
}

func (bot *TelegramRobot) Name() string {
	return "benhuang_bot"
}

func (bot *TelegramRobot) Memo() string {
	return ""
}

func (bot *TelegramRobot) SetChatID(chatID int64) int64 {
	bot.chatID = chatID
	return bot.chatID
}

func (bot *TelegramRobot) SendMessage(respMessage model.Message) error {
	//add time
	bot.setRespMessage(respMessage)
	msg := tgbotapi.NewMessage(bot.chatID, respMessage.String())
	if _, err := bot.Robot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (bot *TelegramRobot) setRespMessage(respMessage model.Message) {
	respMessage.SetTime()

}
