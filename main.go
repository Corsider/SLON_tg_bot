package main

import (
	"SLON_tg_bot/src/app"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	token := os.Getenv("SLON_TOKEN")
	if token == "" {
		panic("empty token")
	}

	bot, err := app.NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Bot.Start(ctx)
}
