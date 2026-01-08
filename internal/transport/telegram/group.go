package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/handlers"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type HandlerDependencies struct {
	FoodCardService *usecases.FoodCardService
	CartService     *usecases.CartService
	StickerService  *usecases.StickerService
	LanguageService *usecases.LanguageService
}

type HandlerGroup struct {
	bot *bot.Bot

	start *handlers.StartHandler
	food  *handlers.FoodCardHandler
	cart  *handlers.CartHandler
	stick *handlers.StickerHandler
	form  *handlers.FoodFormHandler
}

func NewHandlerGroup(
	bot *bot.Bot,
	deps HandlerDependencies,
) *HandlerGroup {
	foodCardHandler := handlers.NewFoodCardHandler(bot, deps.FoodCardService, deps.CartService, deps.LanguageService)
	cartHandler := handlers.NewCartHandler(bot, deps.CartService, deps.LanguageService)
	stickerHandler := handlers.NewStickerHandler(bot, deps.StickerService, deps.LanguageService)
	foodFormHandler := handlers.NewFoodFormHandler(bot, deps.FoodCardService, deps.StickerService, deps.LanguageService)

	cartHandler.SetBackHandler(foodCardHandler.GetOnSelect())
	foodCardHandler.SetCartHandler(cartHandler)
	foodCardHandler.SetAddFoodHandler(foodFormHandler.GetOnAdd())
	foodCardHandler.SetEditFoodHandler(func(card entities.FoodCard) inline.OnSelect {
		return foodFormHandler.GetOnEdit(card.Key)
	})
	foodFormHandler.SetAfterSave(func(ctx context.Context, chatID int64) {
		foodCardHandler.SendFoodMenu(ctx, chatID)
	})
	foodFormHandler.SetAfterCancel(func(ctx context.Context, chatID int64) {
		foodCardHandler.SendFoodMenu(ctx, chatID)
	})

	startKeyboardCfg := keyboards.StartKeyboardConfig{
		FoodHandler:    foodCardHandler.GetOnSelect(),
		StickerHandler: stickerHandler.GetOnSelect(),
	}

	startHandler := handlers.NewStartHandler(bot, startKeyboardCfg, deps.LanguageService)
	startHandler.SetLanguageHandler(startHandler.GetOnSwitchLanguage())
	startKeyboardCfg.LanguageHandler = startHandler.GetOnSwitchLanguage()
	foodCardHandler.SetMainMenuHandler(startHandler.GetOnSelect())
	stickerHandler.SetMainMenuHandler(startHandler.GetOnSelect())

	return &HandlerGroup{
		bot:   bot,
		start: startHandler,
		food:  foodCardHandler,
		cart:  cartHandler,
		stick: stickerHandler,
		form:  foodFormHandler,
	}
}

func (hg *HandlerGroup) Init() {
	hg.start.Init()
	hg.form.Init()
	hg.stick.Init()
}
