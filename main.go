package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"siddharthroy.com/GoodiesProviderBot/internal/proxy"
)

type Application struct {
	bot    *tgbotapi.BotAPI
	client *http.Client
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	client, _, err := proxy.CreateSpysProxyClient()
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	application := Application{
		bot:    botAPI,
		client: client,
	}

	// Set this to true to log all interactions with telegram servers
	botAPI.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := botAPI.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates, application)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel, application Application) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update, application)
		}
	}
}

func handleUpdate(update tgbotapi.Update, application Application) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message, application)

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery, application)
	}
}

func handleMessage(message *tgbotapi.Message, application Application) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var chatId = message.Chat.ID

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(chatId, text, application)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string, application Application) error {
	var err error

	command = strings.Replace(command, application.bot.Self.UserName, "", 1)
	command = strings.Replace(command, "@", "", 1)

	fmt.Println(command)

	switch command {
	case "/menu":
		err = application.HandleMenu(chatId)
	case "/subscribe":
		err = application.HandleSubscribe(chatId)
	case "/unsubcribe":
		err = application.HandleSubscribe(chatId)
	case "/download":
		err = application.HandleDownload(chatId)
	case "/goon":
		err = application.HandleGoon(chatId)
	case "/status":
		err = application.HandleStatus(chatId)
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery, application Application) {

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	application.bot.Send(callbackCfg)

	switch query.Data {
	case SubscribeCommand:
	case UnsubcribeCommand:
	case GoonCommand:
	case StatusCommand:
	case DownloadCommand:
	}

}
