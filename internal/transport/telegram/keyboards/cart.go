package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

func GetCartKeyboard(b *bot.Bot, lang entities.LanguageCode, backHandler, clearCart, acceptCart inline.OnSelect) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button(utils.Text(lang, utils.KeyBtnCartAccept), []byte("accept-cart"), acceptCart).
		Row().
		Button(utils.Text(lang, utils.KeyBtnCartClear), []byte("clear-trash"), clearCart).
		Row().
		Button(utils.Text(lang, utils.KeyBtnBack), []byte("to-food-picking"), backHandler)
	return kb
}

func GetBackKeyboard(b *bot.Bot, lang entities.LanguageCode, backHandler inline.OnSelect) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Button(utils.Text(lang, utils.KeyBtnBack), []byte("to-food-picking"), backHandler)
	return kb
}

func GetOrderDecisionKeyboard(b *bot.Bot, lang entities.LanguageCode, acceptHandler, partialHandler, declineHandler inline.OnSelect) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if acceptHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderAccept), []byte("order-accept"), acceptHandler)
	}
	if partialHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderPartial), []byte("order-partial"), partialHandler)
	}
	if declineHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderDecline), []byte("order-decline"), declineHandler)
	}
	return kb
}

func GetOrderApprovalKeyboard(b *bot.Bot, lang entities.LanguageCode, approveHandler, cancelHandler inline.OnSelect) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if approveHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderApprove), []byte("order-approve"), approveHandler)
	}
	if cancelHandler != nil {
		kb.Row().
			Button(utils.Text(lang, utils.KeyBtnOrderCancel), []byte("order-cancel"), cancelHandler)
	}
	return kb
}
