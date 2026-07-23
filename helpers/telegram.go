package helpers

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegram(msg string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	log.Println("TOKEN =", token)
	log.Println("CHAT ID =", chatIDStr)
	log.Println("PESAN =", msg)

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println("CHAT ID ERROR:", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("BOT ERROR:", err)
		return
	}

	message := tgbotapi.NewMessage(chatID, msg)

	result, err := bot.Send(message)
	if err != nil {
		log.Println("SEND ERROR:", err)
		return
	}

	log.Printf("TELEGRAM BERHASIL: %+v\n", result)
}
