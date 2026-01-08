package keyboards

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

type StartKeyboardConfig struct {
	FoodHandler        inline.OnSelect
	StickerHandler     inline.OnSelect
	AddFoodHandler     inline.OnSelect
	OrderHandler       inline.OnSelect
	MarketplaceHandler inline.OnSelect
	LanguageHandler    inline.OnSelect
}

func GetStartKeyboard(b *bot.Bot, cfg StartKeyboardConfig, lang entities.LanguageCode) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if cfg.FoodHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnFoodMenu), []byte("food_start"), cfg.FoodHandler)
	}
	if cfg.StickerHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnStickers), []byte("stickers_start"), cfg.StickerHandler)
	}
	if cfg.AddFoodHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnAddFood), []byte("food_add"), cfg.AddFoodHandler)
	}
	if cfg.OrderHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderFood), []byte("order_food"), cfg.OrderHandler)
	}
	if cfg.MarketplaceHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnMarketplace), []byte("marketplace_links"), cfg.MarketplaceHandler)
	}
	if cfg.LanguageHandler != nil {
		label := fmt.Sprintf("%s %s", utils.Text(lang, utils.KeyBtnSwitchLanguage), utils.LanguageEmoji(lang))
		kb.Row().
			Button(label, []byte("switch-language"), cfg.LanguageHandler)
	}
	return kb
}
