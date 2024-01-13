package services

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func TelegramBot(config *viper.Viper) {

	if !config.GetBool("TELEGRAM_BOT_ENABLE") {
		log.Println("Telegram bot not enabled")
		return
	}

	bot, err := tgbotapi.NewBotAPI(config.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Println(err)
	}

	bot.Debug = false

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			log.Println(update.Message.From.UserName, update.Message.Chat.ID)
			go data.SetTelegramChat(update.Message.From.UserName, update.Message.Chat.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
