package wa

import (
	"context"
	"fmt"
	"main/ai"
	"main/models"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// variabel untuk client whatsapp
var clientWa *whatsmeow.Client
var DB *gorm.DB

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())

		fmt.Println(" => dari saya =", v.Info.IsFromMe)
		fmt.Println(" => server =", v.Info.MessageSource.Chat.Server)
		fmt.Println(" => apakah group =", v.Info.IsGroup)
		fmt.Println(" => apakah broadcast =", v.Info.IsIncomingBroadcast())

		// filter pesan
		if !v.Info.IsFromMe &&
			v.Info.MessageSource.Chat.Server == "lid" &&
			!v.Info.IsGroup &&
			!v.Info.IsIncomingBroadcast() {

			fmt.Println("PENGIRIM =", v.Info.Sender.User)

			pesan := v.Message.GetConversation()
			fmt.Println("PESAN = " + pesan)

			// membuat array id_wa
			var id_wa []string
			id_wa = append(id_wa, v.Info.ID)

			// status pesan dibaca
			clientWa.MarkRead(
				context.Background(),
				id_wa,
				time.Now(),
				v.Info.Chat,
				v.Info.Sender,
			)

			// pengirim akan menerima status
			clientWa.SubscribePresence(context.Background(), v.Info.Sender)

			// status online
			clientWa.SendPresence(context.Background(), types.PresenceAvailable)

			// jeda 2 detik
			time.Sleep(2 * time.Second)

			// status mengetik
			clientWa.SendChatPresence(
				context.Background(),
				v.Info.Sender,
				types.ChatPresenceComposing,
				types.ChatPresenceMediaText,
			)

			// jeda 3 detik
			time.Sleep(3 * time.Second)

			// status berhenti mengetik
			clientWa.SendChatPresence(
				context.Background(),
				v.Info.Sender,
				types.ChatPresencePaused,
				types.ChatPresenceMediaText,
			)

			pesanAsli := pesan
			pesan = strings.ToLower(pesan)

			if strings.HasPrefix(pesan, "[ai]") {
				pertanyaan := strings.TrimSpace(pesanAsli[4:])

				if pertanyaan != "" {
					jawabanAi := ai.TanyaAi(v.Info.Sender.User, pertanyaan)
					kirimPesanText(v.Info.Sender, jawabanAi)
				} else {
					kirimPesanText(
						v.Info.Sender,
						`Masukkan pertanyaan setelah prefiks [ai].

Contoh: *[ai]Selamat pagi*`,
					)
				}
			} else if pesan == "tes" {
				KirimPesan(v.Info.Sender)
			} else {
				kirimPesanDatabase(v.Info.Sender, pesan)
			}
		}
	}
}

func KoneksiWa(db *gorm.DB) {
	// |------------------------------------------------------------------------------------------------------|
	// | NOTE: You must also import the appropriate DB connector, e.g. github.com/mattn/go-sqlite3 for SQLite |
	// |------------------------------------------------------------------------------------------------------|

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()

	container, err := sqlstore.New(
		ctx,
		"sqlite3",
		"file:examplestore.db?_foreign_keys=on",
		dbLog,
	)
	if err != nil {
		panic(err)
	}

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	if deviceStore != nil {
		deviceStore.Platform = "Windows"
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	clientWa = client
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())

		err = client.Connect()
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	clientWa = client
	DB = db

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func KirimPesan(IDPenerima types.JID) {
	clientWa.SendMessage(
		context.Background(),
		IDPenerima,
		&waE2E.Message{
			Conversation: proto.String("[UJI COBA] \n PESAN OTOMATIS"),
		},
	)
}

func kirimPesanText(idPenerima types.JID, isiPesan string) {
	clientWa.SendMessage(
		context.Background(),
		idPenerima,
		&waE2E.Message{
			Conversation: proto.String(isiPesan),
		},
	)
}

func kirimPesanDatabase(idPenerima types.JID, kode string) {
	var pesan models.Pesan

	result := DB.Where("kode = ?", kode).First(&pesan)
	if result.Error != nil {
		return
	}

	kirimPesanText(idPenerima, pesan.Balasan)
}