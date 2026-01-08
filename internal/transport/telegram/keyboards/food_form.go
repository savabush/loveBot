package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

type FoodFormKeyboardConfig struct {
	SaveHandler         inline.OnSelect
	CancelHandler       inline.OnSelect
	EditNameHandler     inline.OnSelect
	EditDescHandler     inline.OnSelect
	EditPriceHandler    inline.OnSelect
	EditCurrencyHandler inline.OnSelect
	EditTimeHandler     inline.OnSelect
	EditPhotoHandler    inline.OnSelect
	SaveLabel           string
	Lang                entities.LanguageCode
}

func GetFoodFormConfirmKeyboard(b *bot.Bot, cfg FoodFormKeyboardConfig) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())

	saveLabel := cfg.SaveLabel
	if saveLabel == "" {
		saveLabel = utils.Text(cfg.Lang, utils.KeyBtnFoodFormSave)
	}
	if cfg.SaveHandler != nil {
		kb.Row().
			Button(saveLabel, []byte("food-form-save"), cfg.SaveHandler)
	}
	if cfg.CancelHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnFoodFormCancel), []byte("food-form-cancel"), cfg.CancelHandler)
	}

	addEditButton := func(text string, handler inline.OnSelect) {
		if handler == nil {
			return
		}
		kb.Row().
			Button(text, []byte("food-form-edit"), handler)
	}

	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditName), cfg.EditNameHandler)
	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditDesc), cfg.EditDescHandler)
	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditPrice), cfg.EditPriceHandler)
	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditCur), cfg.EditCurrencyHandler)
	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditTime), cfg.EditTimeHandler)
	addEditButton(utils.Text(cfg.Lang, utils.KeyBtnFoodFormEditPhoto), cfg.EditPhotoHandler)

	return kb
}
