package services

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type TelegramNotifier struct {
	ChatID int64
	Config *viper.Viper
}

func NewTelegramNotifier(config *viper.Viper, chatid int64) *TelegramNotifier {
	return &TelegramNotifier{
		ChatID: chatid,
		Config: config,
	}
}

func (e *TelegramNotifier) SendNotification(message string) error {
	botToken := e.Config.GetString("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	bot.Debug = false

	msg := tgbotapi.NewMessage(int64(e.ChatID), message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Ir al sitio web", "https://www.google.com/"),
		),
	)

	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (e *TelegramNotifier) ResetPassword(message string, url string) error {

	if e.ChatID == 0 {
		return errors.New("no existe el chat id")
	}

	botToken := e.Config.GetString("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	bot.Debug = false

	msg := tgbotapi.NewMessage(e.ChatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	// boton para cambiar la contraseña
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Cambiar contraseña", url),
		),
	)

	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
