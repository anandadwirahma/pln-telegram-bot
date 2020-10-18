package main

import (
	"log"

	"github.com/yanzay/tbot"
)

type application struct {
	client *tbot.Client
}

const (
	TELEGRAM_TOKEN = "1135438873:AAHOKp7H_-VYRI8n6QHipyLEjz-v52j6pjc"
	BOX_URL        = "https://boxofimagination.com"
	//WEBHHOK_URL    = "https://e8fb0cfb8445.ngrok.io"
	//ADDR           = "0.0.0.0:8000"
)

var (
	app application
	bot *tbot.Server
)

func main() {
	//webhook := tbot.WithWebhook(WEBHHOK_URL, ADDR)
	bot = tbot.New(TELEGRAM_TOKEN)
	app.client = bot.Client()

	bot.HandleMessage(".*halo.*", app.greatingHandler)
	bot.HandleMessage(".*.*", app.commonHandler)

	bot.HandleCallback(app.callbackHandler)
	log.Fatal(bot.Start())
}
