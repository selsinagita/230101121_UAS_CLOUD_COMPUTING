package helpers

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegram(msg string) {

	token := os.Getenv("TELEGRAM_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	fmt.Println("TOKEN :", token)
	fmt.Println("CHAT ID :", chatIDStr)

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("CHAT ID ERROR :", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("BOT ERROR :", err)
		return
	}

	message := tgbotapi.NewMessage(
		chatID,
		msg,
	)

	_, err = bot.Send(message)
	if err != nil {
		fmt.Println("SEND ERROR :", err)
		return
	}

	fmt.Println("TELEGRAM BERHASIL DIKIRIM")
}
