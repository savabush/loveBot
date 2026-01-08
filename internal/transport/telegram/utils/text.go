package utils

import "github.com/savabush/breakfastLoveBot/internal/entities"

type TextKey string

const (
	KeyDefaultMessageTextWithStickers TextKey = "default_with_stickers"
	KeyDefaultMessageText             TextKey = "default"
	KeyStartMessageText               TextKey = "start_message"
	KeyFoodCardText                   TextKey = "food_card"
	KeyEmptyFoodCardText              TextKey = "food_empty"
	KeyFoodAddedToCartText            TextKey = "food_added"
	KeyFoodDeletedText                TextKey = "food_deleted"
	KeyUnknownFoodSelectedText        TextKey = "food_unknown"
	KeyEmptyCartText                  TextKey = "cart_empty"
	KeyCartListText                   TextKey = "cart_list"
	KeyCartAcceptedForOtherText       TextKey = "cart_accepted_other"
	KeyCartOrderAcceptedText          TextKey = "cart_order_accepted"
	KeyCartOrderDeclinedText          TextKey = "cart_order_declined"
	KeyCartPartialSelectionText       TextKey = "cart_partial_select"
	KeyCartPartialSendEmptyText       TextKey = "cart_partial_empty"
	KeyCartPartialSentText            TextKey = "cart_partial_sent"
	KeyCartPartialRequestText         TextKey = "cart_partial_request"
	KeyCartPartialApprovedText        TextKey = "cart_partial_approved"
	KeyCartPartialDeclinedText        TextKey = "cart_partial_declined"
	KeyCartDecisionSentText           TextKey = "cart_decision_sent"
	KeyCartAcceptedText               TextKey = "cart_accepted"
	KeyCartClearedText                TextKey = "cart_cleared"
	KeyCartItemDetails                TextKey = "cart_item_details"
	KeyEmptyStickersText              TextKey = "stickers_empty"
	KeyFoodFormStartAddText           TextKey = "food_form_start"
	KeyFoodFormNamePrompt             TextKey = "food_form_name"
	KeyFoodFormDescriptionPrompt      TextKey = "food_form_desc"
	KeyFoodFormPricePrompt            TextKey = "food_form_price"
	KeyFoodFormCurrencyPrompt         TextKey = "food_form_currency"
	KeyFoodFormTimePrompt             TextKey = "food_form_time"
	KeyFoodFormPhotoPrompt            TextKey = "food_form_photo"
	KeyFoodFormSummaryText            TextKey = "food_form_summary"
	KeyFoodFormCardNotFoundText       TextKey = "food_form_not_found"
	KeyFoodFormUnknownInputText       TextKey = "food_form_unknown"
	KeyFoodFormSavedText              TextKey = "food_form_saved"
	KeyFoodFormUpdatedText            TextKey = "food_form_updated"
	KeyFoodFormCancelledText          TextKey = "food_form_cancelled"
	KeyStickerFormPrompt              TextKey = "sticker_form_prompt"
	KeyStickerFormSavedText           TextKey = "sticker_form_saved"
	KeyStickerDeletedText             TextKey = "sticker_deleted"
	KeyStickerDeleteEmpty             TextKey = "sticker_delete_empty"
	KeyStickerCatchText               TextKey = "sticker_catch"
	KeyCartAcceptedInfoText           TextKey = "cart_accept_info"

	KeyBtnFoodMenu          TextKey = "btn_food_menu"
	KeyBtnStickers          TextKey = "btn_stickers"
	KeyBtnAddFood           TextKey = "btn_add_food"
	KeyBtnOrderFood         TextKey = "btn_order_food"
	KeyBtnMarketplace       TextKey = "btn_marketplace"
	KeyBtnSwitchLanguage    TextKey = "btn_switch_language"
	KeyBtnRandomFood        TextKey = "btn_random_food"
	KeyBtnAddToCart         TextKey = "btn_add_to_cart"
	KeyBtnEditFood          TextKey = "btn_edit_food"
	KeyBtnDeleteFood        TextKey = "btn_delete_food"
	KeyBtnShowCart          TextKey = "btn_show_cart"
	KeyBtnAddFoodEmpty      TextKey = "btn_add_food_empty"
	KeyBtnMainMenu          TextKey = "btn_main_menu"
	KeyBtnBack              TextKey = "btn_back"
	KeyBtnCartAccept        TextKey = "btn_cart_accept"
	KeyBtnCartClear         TextKey = "btn_cart_clear"
	KeyBtnNextSticker       TextKey = "btn_next_sticker"
	KeyBtnAddSticker        TextKey = "btn_add_sticker"
	KeyBtnDeleteSticker     TextKey = "btn_delete_sticker"
	KeyBtnOrderAccept       TextKey = "btn_order_accept"
	KeyBtnOrderPartial      TextKey = "btn_order_partial"
	KeyBtnOrderDecline      TextKey = "btn_order_decline"
	KeyBtnOrderApprove      TextKey = "btn_order_approve"
	KeyBtnOrderCancel       TextKey = "btn_order_cancel"
	KeyBtnPartialSend       TextKey = "btn_partial_send"
	KeyBtnPartialCancel     TextKey = "btn_partial_cancel"
	KeyBtnFoodFormSave      TextKey = "btn_food_form_save"
	KeyBtnFoodFormUpdate    TextKey = "btn_food_form_update"
	KeyBtnFoodFormCancel    TextKey = "btn_food_form_cancel"
	KeyBtnFoodFormEditName  TextKey = "btn_food_form_edit_name"
	KeyBtnFoodFormEditDesc  TextKey = "btn_food_form_edit_desc"
	KeyBtnFoodFormEditPrice TextKey = "btn_food_form_edit_price"
	KeyBtnFoodFormEditCur   TextKey = "btn_food_form_edit_cur"
	KeyBtnFoodFormEditTime  TextKey = "btn_food_form_edit_time"
	KeyBtnFoodFormEditPhoto TextKey = "btn_food_form_edit_photo"
)

var translations = map[entities.LanguageCode]map[TextKey]string{
	entities.LanguageEN: {
		KeyDefaultMessageTextWithStickers: `
Please, type ||\(I'm begging\)|| */start* 😇

And here are some cat stickers 👇
`,
		KeyDefaultMessageText: `
Please, type ||\(I'm begging\)|| */start* 😇

There could be some cat stickers 😥
||You can add them with the */start* buttons||
`,
		KeyStartMessageText: `
*%v, %v\!* ♥ 

Today I suggest choosing what you want to eat\! 🍰

The menu is still __growing__, because our clever _developer_ is thinking what else to add 😏

||Secretly: we will also move family wishlist items and shopping things here soon\!\! 💅
There is much more waiting for you 🦍||
	`,
		KeyFoodCardText: `
*%v*

%v 

%v %v
%v minutes
	`,
		KeyEmptyFoodCardText: `
Sadly, there is no food yet.
Let's add our first dishes\!
	`,
		KeyFoodAddedToCartText:     `Added to cart ✅`,
		KeyFoodDeletedText:         `Dish deleted 🗑️`,
		KeyUnknownFoodSelectedText: `Couldn't find the dish, so I showed the first one`,
		KeyEmptyCartText: `
The cart is empty 🥲
Add something tasty\!
	`,
		KeyCartListText: `
In the cart now:
%v
	`,
		KeyCartAcceptedForOtherText: `From your half ❤️
Need to make:
%v
Cooking time: %v minutes
Your half owes you:
%v`,
		KeyCartOrderAcceptedText: `Your half accepted the order ✅
Need to make:
%v
Cooking time: %v minutes
Your half owes you:
%v`,
		KeyCartOrderDeclinedText:    `Your half declined ❌`,
		KeyCartPartialSelectionText: `Pick the dishes you will cook:`,
		KeyCartPartialSendEmptyText: `Nothing selected ⚠️`,
		KeyCartPartialSentText:      `Sent for approval ✅`,
		KeyCartPartialRequestText: `Your half suggests changes:
Need to make:
%v
Cooking time: %v minutes
Your half owes you:
%v`,
		KeyCartPartialApprovedText: `Order confirmed ✅
Need to make:
%v
Cooking time: %v minutes
Your half owes you:
%v`,
		KeyCartPartialDeclinedText: `No, I don't want this order ❌`,
		KeyCartDecisionSentText:    `Response sent ✅`,
		KeyCartAcceptedText:        `Order accepted ✅`,
		KeyCartClearedText:         `Cart cleared 🧹`,
		KeyCartItemDetails:         `Price: %v, Time: %d min`,
		KeyEmptyStickersText: `
I don't have stickers yet 🥺
Please add IDs to data/stickers\.json
	`,
		KeyFoodFormStartAddText: `
Let's add a new dish 🤌
I'll ask one field at a time, just reply 😇
Type */cancel* if you change your mind
	`,
		KeyFoodFormNamePrompt:        `What should we call the dish?`,
		KeyFoodFormDescriptionPrompt: `Describe it tasty 🍝`,
		KeyFoodFormPricePrompt:       `Enter the price \(you can give multiple numbers separated by spaces\)`,
		KeyFoodFormCurrencyPrompt:    `Which currency? For example, ₽ or RUB`,
		KeyFoodFormTimePrompt:        `How many minutes to cook?`,
		KeyFoodFormPhotoPrompt:       `Send a photo of the dish. If none — type "\-"`,
		KeyFoodFormSummaryText: `
*Let's check everything:*

*Name:* %v
*Description:* %v
*Price:* %v
*Currency:* %v
*Cooking time:* %v minutes
*Photo:* %v
	`,
		KeyFoodFormCardNotFoundText: `Dish not found, try again 🙈`,
		KeyFoodFormUnknownInputText: `I'm waiting for the dish form input. If you want to exit — type /cancel`,
		KeyFoodFormSavedText:        `New dish saved\! 🥳`,
		KeyFoodFormUpdatedText:      `Dish updated ✅`,
		KeyFoodFormCancelledText:    `Ok, nothing saved`,
		KeyStickerFormPrompt:        `Send a sticker, and I'll add it to the collection`,
		KeyStickerFormSavedText:     `Sticker added\!`,
		KeyStickerDeletedText:       `Sticker deleted 🗑️`,
		KeyStickerDeleteEmpty:       `No sticker to delete`,
		KeyStickerCatchText:         `Here is a kitty 🐾`,

		KeyBtnFoodMenu:          `Yummies 🥞`,
		KeyBtnStickers:          `Stickers 🐾`,
		KeyBtnAddFood:           `Add dish ➕`,
		KeyBtnOrderFood:         `Want to order food 📱`,
		KeyBtnMarketplace:       `Marketplace monsters 👾`,
		KeyBtnSwitchLanguage:    `Сменить язык / Switch language`,
		KeyBtnRandomFood:        `🎲 Random dish`,
		KeyBtnAddToCart:         `➕ To cart`,
		KeyBtnEditFood:          `✏️ Edit dish`,
		KeyBtnDeleteFood:        `🗑️ Delete dish`,
		KeyBtnShowCart:          `Cart 🛒`,
		KeyBtnAddFoodEmpty:      `Add dish ➕`,
		KeyBtnMainMenu:          `🏠 Main menu`,
		KeyBtnBack:              `Back`,
		KeyBtnCartAccept:        `Place order ✅`,
		KeyBtnCartClear:         `Clear cart 🗑️`,
		KeyBtnNextSticker:       `Another 🐾`,
		KeyBtnAddSticker:        `Add sticker ➕`,
		KeyBtnDeleteSticker:     `🗑️ Delete sticker`,
		KeyBtnOrderAccept:       `✅ I'll do it`,
		KeyBtnOrderPartial:      `🟡 Partially`,
		KeyBtnOrderDecline:      `❌ Decline`,
		KeyBtnOrderApprove:      `✅ Confirm changes`,
		KeyBtnOrderCancel:       `❌ Cancel order`,
		KeyBtnPartialSend:       `✅ Send for approval`,
		KeyBtnPartialCancel:     `↩️ Cancel`,
		KeyBtnFoodFormSave:      `Save ✅`,
		KeyBtnFoodFormUpdate:    `Update ✅`,
		KeyBtnFoodFormCancel:    `Cancel ❌`,
		KeyBtnFoodFormEditName:  `✏️ Edit name`,
		KeyBtnFoodFormEditDesc:  `✏️ Edit description`,
		KeyBtnFoodFormEditPrice: `✏️ Edit price`,
		KeyBtnFoodFormEditCur:   `✏️ Edit currency`,
		KeyBtnFoodFormEditTime:  `✏️ Edit cooking time`,
		KeyBtnFoodFormEditPhoto: `✏️ Edit photo`,
	},
	entities.LanguageRU: {
		KeyDefaultMessageTextWithStickers: `
Напиши, пожалуйста, ||\(молю\)|| */start* 😇

Вот, кстати, стикеры с котиками 👇
`,
		KeyDefaultMessageText: `
Напиши, пожалуйста, ||\(молю\)|| */start* 😇

Тут могли быть стикеры с котиками 😥
||Их можно добавить с помощью кнопочек при написании команды */start*||
`,
		KeyStartMessageText: `
*%v, %v\!* ♥ 

Сегодня я предлагаю тебе выбрать, что ты хочешь кушац\! 🍰

Выбор блюд еще __расширяется__, так как наш умный _разработчик_ думает, что можно крутого во мне сделать\! 😏

||По секрету скажу, что сюда переедут семейные хотелки и вещи с озона и других маркетплейсов в ближайшем будущем\!\! 💅
Многое другое еще тебя ждет тут 🦍||
	`,
		KeyFoodCardText: `
*%v*

%v 

%v %v
%v минутов
	`,
		KeyEmptyFoodCardText: `
Увы, но ничего не добавлено по еде
Давай добавим наши первые блюда\!
	`,
		KeyFoodAddedToCartText:     `Блюдо добавлено в корзину ✅`,
		KeyFoodDeletedText:         `Блюдо удалено 🗑️`,
		KeyUnknownFoodSelectedText: `Не смог найти блюдо, поэтому показал первое из списка`,
		KeyEmptyCartText: `
В корзине пока пусто 🥲
Добавь что‑нибудь вкусное\!
	`,
		KeyCartListText: `
В корзине сейчас:
%v
	`,
		KeyCartAcceptedForOtherText: `От половинки ❤️
Нужно сделать:
%v
Время готовки: %v минут
Половинка должна тебе за это:
%v`,
		KeyCartOrderAcceptedText: `Половинка взяла заказ ✅
Нужно сделать:
%v
Время готовки: %v минут
Половинка должна тебе за это:
%v`,
		KeyCartOrderDeclinedText:    `Половинка отказалась ❌`,
		KeyCartPartialSelectionText: `Выбери блюда, которые готовишь:`,
		KeyCartPartialSendEmptyText: `Ничего не выбрано ⚠️`,
		KeyCartPartialSentText:      `Отправила на аппрув ✅`,
		KeyCartPartialRequestText: `Половинка предлагает изменения:
Нужно сделать:
%v
Время готовки: %v минут
Половинка должна тебе за это:
%v`,
		KeyCartPartialApprovedText: `Заказ подтвержден ✅
Нужно сделать:
%v
Время готовки: %v минут
Половинка должна тебе за это:
%v`,
		KeyCartPartialDeclinedText: `Нет, заказ такой не хочется ❌`,
		KeyCartDecisionSentText:    `Ответ отправлен ✅`,
		KeyCartAcceptedText:        `Заказ принят ✅`,
		KeyCartClearedText:         `Корзина очищена 🧹`,
		KeyCartItemDetails:         `Цена: %v, Время: %d мин`,
		KeyEmptyStickersText: `
Пока что у меня нет стикеров 🥺
Добавь айдишники в data/stickers\.json, пожалуйста
	`,
		KeyFoodFormStartAddText: `
Давай добавим новое блюдо 🤌
Я буду спрашивать по одному полю, а ты просто отвечай 😇
Напиши */cancel*, если передумаешь
	`,
		KeyFoodFormNamePrompt:        `Как назовем блюдо?`,
		KeyFoodFormDescriptionPrompt: `Опиши блюдо вкусно 🍝`,
		KeyFoodFormPricePrompt:       `Введи цену \(можно несколько чисел через пробел\)`,
		KeyFoodFormCurrencyPrompt:    `В какой валюте? Например, ₽ или RUB`,
		KeyFoodFormTimePrompt:        `Сколько минут готовится блюдо?`,
		KeyFoodFormPhotoPrompt:       `Пришли фото блюда \(картинку\)\. Если нет — напиши "\-"`,
		KeyFoodFormSummaryText: `
*Проверим, все ли верно?*

*Название:* %v
*Описание:* %v
*Цена:* %v
*Валюта:* %v
*Время готовки:* %v минут
*Фото:* %v
	`,
		KeyFoodFormCardNotFoundText: `Не нашел блюдо, попробуй еще раз 🙈`,
		KeyFoodFormUnknownInputText: `Я жду ответ на вопрос про блюдо. Если хочешь выйти — напиши /cancel`,
		KeyFoodFormSavedText:        `Новое блюдо сохранено\! 🥳`,
		KeyFoodFormUpdatedText:      `Блюдо обновлено ✅`,
		KeyFoodFormCancelledText:    `Окей, ничего не сохраняю`,
		KeyStickerFormPrompt:        `Пришли стикер, и я добавлю его в коллекцию`,
		KeyStickerFormSavedText:     `Стикер добавлен\!`,
		KeyStickerDeletedText:       `Стикер удален 🗑️`,
		KeyStickerDeleteEmpty:       `Нет стикера для удаления`,
		KeyStickerCatchText:         `Лови котейку 🐾`,

		KeyBtnFoodMenu:          `Вкусняшки 🥞`,
		KeyBtnStickers:          `Стикеры 🐾`,
		KeyBtnAddFood:           `Добавить блюдо ➕`,
		KeyBtnOrderFood:         `Хочется заказать еды 📱`,
		KeyBtnMarketplace:       `Маркетплейс монстры 👾`,
		KeyBtnSwitchLanguage:    `Сменить язык / Switch language`,
		KeyBtnRandomFood:        `🎲 Случайное блюдо`,
		KeyBtnAddToCart:         `➕ В корзину`,
		KeyBtnEditFood:          `✏️ Редактировать блюдо`,
		KeyBtnDeleteFood:        `🗑️ Удалить блюдо`,
		KeyBtnShowCart:          `Корзина 🛒`,
		KeyBtnAddFoodEmpty:      `Добавить блюдо ➕`,
		KeyBtnMainMenu:          `🏠 В главное меню`,
		KeyBtnBack:              `Назад`,
		KeyBtnCartAccept:        `Оформить заказ ✅`,
		KeyBtnCartClear:         `Очистить корзину 🗑️`,
		KeyBtnNextSticker:       `Еще один 🐾`,
		KeyBtnAddSticker:        `Добавить стикер ➕`,
		KeyBtnDeleteSticker:     `🗑️ Удалить стикер`,
		KeyBtnOrderAccept:       `✅ Хорошо, сделаю`,
		KeyBtnOrderPartial:      `🟡 Частично`,
		KeyBtnOrderDecline:      `❌ Отказ`,
		KeyBtnOrderApprove:      `✅ Подтвердить изменения`,
		KeyBtnOrderCancel:       `❌ Аннулирование заказа`,
		KeyBtnPartialSend:       `✅ Отправить на аппрув`,
		KeyBtnPartialCancel:     `↩️ Отмена`,
		KeyBtnFoodFormSave:      `Сохранить ✅`,
		KeyBtnFoodFormUpdate:    `Обновить ✅`,
		KeyBtnFoodFormCancel:    `Отменить ❌`,
		KeyBtnFoodFormEditName:  `✏️ Изменить название`,
		KeyBtnFoodFormEditDesc:  `✏️ Изменить описание`,
		KeyBtnFoodFormEditPrice: `✏️ Изменить цену`,
		KeyBtnFoodFormEditCur:   `✏️ Изменить валюту`,
		KeyBtnFoodFormEditTime:  `✏️ Изменить время готовки`,
		KeyBtnFoodFormEditPhoto: `✏️ Изменить фото`,
	},
}

func Text(lang entities.LanguageCode, key TextKey) string {
	if lang == "" {
		lang = entities.LanguageEN
	}
	if m, ok := translations[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if v, ok := translations[entities.LanguageEN][key]; ok {
		return v
	}
	return ""
}

func LanguageEmoji(lang entities.LanguageCode) string {
	switch lang {
	case entities.LanguageRU:
		return "🇷🇺"
	default:
		return "🇬🇧"
	}
}
