package ai

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func MulaiChatAi() {
	client := openai.NewClient(
		option.WithBaseURL("http://localhost:8080/v1/"),
		option.WithAPIKey("saksake-karena-gak-butuh-api-key"),
	)

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	// Menyimpan riwayat percakapan
	var chatHistory []openai.ChatCompletionMessageParamUnion

	for {
		// Input dari pengguna
		fmt.Print("\nAnda: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Keluar dari chatbot
		if strings.EqualFold(userInput, "keluar") {
			fmt.Println("Terima kasih telah menggunakan FikomBot.")
			break
		}

		// Tambahkan pesan pengguna ke histori
		chatHistory = append(chatHistory, openai.UserMessage(userInput))

		fmt.Print("FikomBot: ")

		// Request streaming ke model AI
		stream := client.Chat.Completions.NewStreaming(
			ctx,
			openai.ChatCompletionNewParams{
				Model:    openai.ChatModel("diisi-saksake"),
				Messages: chatHistory,
			},
		)

		var fullResponse strings.Builder

		// Membaca token hasil streaming
		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) > 0 {
				token := chunk.Choices[0].Delta.Content
				fmt.Print(token)
				fullResponse.WriteString(token)
			}
		}

		// Cek apakah terjadi error saat streaming
		if err := stream.Err(); err != nil {
			fmt.Printf("\n[Error saat streaming: %v]\n", err)
			continue
		}

		// Simpan jawaban AI ke histori agar percakapan tetap berlanjut
		chatHistory = append(
			chatHistory,
			openai.AssistantMessage(fullResponse.String()),
		)

		fmt.Println()
	}
}