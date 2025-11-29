package breakfastLoveBot

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/savabush/breakfastLoveBot/internal/config"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/handlers"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/middlewares"
)

func App() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(
			middlewares.LogMessage,
			middlewares.CheckUserIDs,
			middlewares.DeletePreviousBotMessage,
		),
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	b, err := bot.New(config.Cfg.TelegramBotToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	hg := telegram.NewHandlerGroup(b)
	hg.Init()

	log.Println("Bot is running...")

	b.Start(ctx)
}
