package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func GetStartKeyboard(b *bot.Bot, foodCardHandler inline.OnSelect) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Вкусняшки 🥞", []byte("food_start"), foodCardHandler)
	// .Row()
	// Button("Хочется заказать еды 📱", b, bot.MatchTypeExact, OrderFood)
	// .Row()
	// Button("Маркетплейс монстры 👾", b, bot.MatchTypeExact, MarketPlace)
	// .Row()
	// Button("Отправить случайный милый стикер половинке", b, bot.MatchTypeExact, RandomStickerToOther)
	// .Row()
	// Button("Идеи для времопровождения", b, bot.MatchTypeExact, RandomStickerToOther)
	// .Row()
	// Button("Рандомайзер", b, bot.MatchTypeExact, RandomPicker)

	return kb
}
