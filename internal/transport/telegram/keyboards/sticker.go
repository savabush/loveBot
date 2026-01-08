package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

type StickerKeyboardConfig struct {
	NextHandler     inline.OnSelect
	AddHandler      inline.OnSelect
	DeleteHandler   inline.OnSelect
	MainMenuHandler inline.OnSelect
	Lang            entities.LanguageCode
}

func GetStickerKeyboard(b *bot.Bot, cfg StickerKeyboardConfig) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if cfg.NextHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnNextSticker), []byte("stickers-next"), cfg.NextHandler)
	}
	if cfg.AddHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnAddSticker), []byte("stickers-add"), cfg.AddHandler)
	}
	if cfg.DeleteHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnDeleteSticker), []byte("stickers-delete"), cfg.DeleteHandler)
	}
	if cfg.MainMenuHandler != nil {
		kb.Row().
			Button(utils.Text(cfg.Lang, utils.KeyBtnMainMenu), []byte("stickers-main-menu"), cfg.MainMenuHandler)
	}
	return kb
}
