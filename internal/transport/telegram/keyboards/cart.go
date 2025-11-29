package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func GetCartKeyboard(b *bot.Bot, backHandler, clearCart, acceptCart inline.OnSelect) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Оформить заказ ✅", []byte("accept-cart"), acceptCart).
		Row().
		Button("Очистить корзину 🗑️", []byte("clear-trash"), clearCart).
		Row().
		Button("Назад", []byte("to-food-picking"), backHandler)
	return kb
}

func GetBackKeyboard(b *bot.Bot, backHandler inline.OnSelect) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Button("Назад", []byte("to-food-picking"), backHandler)
	return kb
}
