package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

type FoodButton struct {
	Text    string
	OnClick inline.OnSelect
}

type FoodKeyboardConfig struct {
	RandomHandler        inline.OnSelect
	AddToCart            inline.OnSelect
	ShowCartHandler      inline.OnSelect
	MainMenuHandler      inline.OnSelect
	DeleteCurrentHandler inline.OnSelect
	CartButtonText       string
	AddFoodHandler       inline.OnSelect
	EditCurrentHandler   inline.OnSelect
	FoodButtons          []FoodButton
}

type EmptyFoodKeyboardConfig struct {
	AddFoodHandler inline.OnSelect
	Lang           entities.LanguageCode
}

func GetFoodKeyboard(b *bot.Bot, cfg FoodKeyboardConfig, lang entities.LanguageCode) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())

	if cfg.RandomHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnRandomFood), []byte("food-random"), cfg.RandomHandler)
	}

	if cfg.AddToCart != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnAddToCart), []byte("food-add-to-cart"), cfg.AddToCart)
	}

	if cfg.EditCurrentHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnEditFood), []byte("food-edit-current"), cfg.EditCurrentHandler)
	}

	if cfg.DeleteCurrentHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnDeleteFood), []byte("food-delete-current"), cfg.DeleteCurrentHandler)
	}

	if cfg.ShowCartHandler != nil {
		text := cfg.CartButtonText
		if text == "" {
			text = utils.Text(lang, utils.KeyBtnShowCart)
		}
		kb.Row().
			Button(text, []byte("food-show-cart"), cfg.ShowCartHandler)
	}

	if cfg.MainMenuHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnMainMenu), []byte("food-main-menu"), cfg.MainMenuHandler)
	}

	for i, btn := range cfg.FoodButtons {
		if i%2 == 0 {
			kb.Row()
		}
		kb.Button(btn.Text, []byte("food-select"), btn.OnClick)
	}

	if cfg.AddFoodHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnAddFood), []byte("food-add"), cfg.AddFoodHandler)
	}

	return kb
}

func GetEmptyFoodKeyboard(b *bot.Bot, cfg EmptyFoodKeyboardConfig) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if cfg.AddFoodHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnAddFoodEmpty), []byte("food-empty-add"), cfg.AddFoodHandler)
	}
	return kb
}
