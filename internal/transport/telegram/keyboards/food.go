package keyboards

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func GetFoodKeyboard(
	b *bot.Bot,
	// addToCart,
	// cancelHandler,
	// showFoodCard,
	// showCart inline.OnSelect,
	countOfCardsInCart uint8,
) *inline.Keyboard {

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row()
		// Button(fmt.Sprintf("Корзина 🛒 (%d)", countOfCardsInCart), []byte("show-cart"), showCart). //.Button("👆 Добавить в корзину", []byte("add-to-cart"), addToCart).
		// Row().
		// 	Button("Бутербродики 🥪", []byte("food-sandwich"), showFoodCard).Button("Скрррембл 🍳", []byte("food-scramble"), showFoodCard).
		// 	Row().
		// 	Button("Пастя лосось 🍝", []byte("food-pasta-losos"), showFoodCard).Button("Лазаньй 🤌🏻", []byte("food-lasagna"), showFoodCard).
		// 	Row().
		// 	Button("Картошка тефтельная 🥔🍖", []byte("food-potato-teftelya"), showFoodCard).Button("Пирог лимонный 🍋", []byte("food-cake-lemon"), showFoodCard).
		// 	Row().
		// 	Button("Вишневи пирог 🍓", []byte("food-cake-strawberry"), showFoodCard).Button("Жаркое соусное 🍖", []byte("food-jarko"), showFoodCard).
		// 	Row().
		// 	Button("Котлеты отрубили 🔪", []byte("food-kotletos"), showFoodCard).Button("Яйцо Павел 🥚", []byte("food-egg-pavel"), showFoodCard).
		// 	Row().
		// 	Button("Пирожучки 🥖", []byte("food-piroshki"), showFoodCard).Button("Петушиный суп 🥘", []byte("food-egg-soup"), showFoodCard).
		// 	Row().
		// 	Button("Запуканочка 🥵", []byte("food-zapekano4ka"), showFoodCard).Button("Соси сосочки 🌭", []byte("food-sosisos"), showFoodCard).
		// 	Row().
		// 	Button("Картошка освобожденная 🍟", []byte("food-fri"), showFoodCard).
		// 	Row().
		// Button("В главное меню", []byte("cancel"), cancelHandler)
	return kb
}
