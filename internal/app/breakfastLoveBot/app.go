package breakfastLoveBot

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/savabush/breakfastLoveBot/internal/config"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/cart"
	"github.com/savabush/breakfastLoveBot/internal/repository/foodCard"
	"github.com/savabush/breakfastLoveBot/internal/repository/language"
	"github.com/savabush/breakfastLoveBot/internal/repository/sticker"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/handlers"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/middlewares"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
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
	}

	foodCardRepo := foodCard.NewMemoryRepository()
	cartRepo := cart.NewMemoryRepository([]entities.UserTelegramID{
		entities.UserTelegramID(config.Cfg.UserID1),
		entities.UserTelegramID(config.Cfg.UserID2),
	})
	langRepo := language.NewMemoryRepository()
	stickerRepo := sticker.NewMemoryRepository()

	deps := telegram.HandlerDependencies{
		FoodCardService: usecases.NewFoodCardService(foodCardRepo),
		CartService:     usecases.NewCartService(cartRepo),
		StickerService:  usecases.NewStickerService(stickerRepo),
		LanguageService: usecases.NewLanguageService(langRepo),
	}

	opts = append(opts, bot.WithDefaultHandler(handlers.NewDefaultHandler(deps.LanguageService)))

	b, err := bot.New(config.Cfg.TelegramBotToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	hg := telegram.NewHandlerGroup(b, deps)
	hg.Init()

	log.Println("Bot is running...")

	b.Start(ctx)
}
