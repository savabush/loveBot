package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/handlers"
)

type HandlerGroup struct {
	bot *bot.Bot

	start *handlers.StartHandler
	// stickers       *StickerHandler
	// foodCard       *FoodCardHandler
	// cart           *CartHandler
}

func NewHandlerGroup(
	bot *bot.Bot,
	// stickerUC *usecases.StickerService,
	// foodCardUC *usecases.FoodCardService,
	// cartUC *usecases.CartService,
) *HandlerGroup {
	return &HandlerGroup{
		bot:   bot,
		start: handlers.NewStartHandler(bot),
		// stickers:       NewStickerHandler(bot, stickerUC),
		// foodCard:       NewFoodCardHandler(bot, foodCardUC),
		// cart:           NewCartHandler(bot, cartUC),
	}
}

func (hg *HandlerGroup) Init() {
	hg.start.Init()
	// hg.stickers.Init()
	// hg.foodCard.Init()
	// hg.cart.Init()
}
