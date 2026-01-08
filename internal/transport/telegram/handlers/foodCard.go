package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type FoodCardHandler struct {
	bot *bot.Bot

	foodCardService *usecases.FoodCardService
	cartService     *usecases.CartService
	langService     *usecases.LanguageService

	cartHandler   *CartHandler
	addFoodButton inline.OnSelect
	editFoodFunc  func(card entities.FoodCard) inline.OnSelect
	mainMenuFn    inline.OnSelect
}

func NewFoodCardHandler(bot *bot.Bot, foodCardService *usecases.FoodCardService, cartService *usecases.CartService, langService *usecases.LanguageService) *FoodCardHandler {
	return &FoodCardHandler{
		bot:             bot,
		foodCardService: foodCardService,
		cartService:     cartService,
		langService:     langService,
	}
}

func (fch *FoodCardHandler) SetCartHandler(cartHandler *CartHandler) {
	fch.cartHandler = cartHandler
}

func (fch *FoodCardHandler) SetAddFoodHandler(handler inline.OnSelect) {
	fch.addFoodButton = handler
}

func (fch *FoodCardHandler) SetEditFoodHandler(handler func(card entities.FoodCard) inline.OnSelect) {
	fch.editFoodFunc = handler
}

func (fch *FoodCardHandler) SetMainMenuHandler(handler inline.OnSelect) {
	fch.mainMenuFn = handler
}

func (fch *FoodCardHandler) GetOnSelect() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		fch.showRandomFood(ctx, bot, mes)
	}
}

func (fch *FoodCardHandler) GetOnSelectFood(key entities.FoodCardKey) inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		fch.showFoodByKey(ctx, bot, mes, key)
	}
}

func (fch *FoodCardHandler) GetOnSelectRandom() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		fch.showRandomFood(ctx, bot, mes)
	}
}

func (fch *FoodCardHandler) SendFoodMenu(ctx context.Context, chatID int64) {
	defer fch.recoverPanic()
	foodCards := fch.foodCardService.GetAllFoodCards()
	lang := fch.langService.Get(entities.UserTelegramID(chatID))
	if len(foodCards) == 0 {
		kb := keyboards.GetEmptyFoodKeyboard(fch.bot, keyboards.EmptyFoodKeyboardConfig{
			AddFoodHandler: fch.addFoodButton,
			Lang:           lang,
		})
		m := utils.NewMessage(ctx, fch.bot, nil)
		m.Send(chatID, utils.Text(lang, utils.KeyEmptyFoodCardText), true, kb)
		return
	}
	randomFoodCard := foodCards[rand.IntN(len(foodCards))]
	msgText := fch.buildFoodText(lang, randomFoodCard, "")
	kb := fch.buildFoodKeyboard(chatID, lang, randomFoodCard, foodCards)
	fch.sendFoodCardMessage(ctx, fch.bot, chatID, randomFoodCard, msgText, kb)
}

func (fch *FoodCardHandler) showRandomFood(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	defer fch.recoverPanic()
	foodCards := fch.foodCardService.GetAllFoodCards()
	lang := fch.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
	if len(foodCards) == 0 {
		fch.renderEmptyState(ctx, b, mes)
		return
	}
	randomFoodCard := foodCards[rand.IntN(len(foodCards))]
	fch.renderFoodCard(ctx, b, mes, randomFoodCard, foodCards, "", lang)
}

func (fch *FoodCardHandler) showFoodByKey(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, key entities.FoodCardKey) {
	defer fch.recoverPanic()
	foodCards := fch.foodCardService.GetAllFoodCards()
	lang := fch.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
	if len(foodCards) == 0 {
		fch.renderEmptyState(ctx, b, mes)
		return
	}
	card, ok := fch.findFoodCard(foodCards, key)
	if !ok {
		fch.renderFoodCard(ctx, b, mes, foodCards[0], foodCards, utils.Text(lang, utils.KeyUnknownFoodSelectedText), lang)
		return
	}
	fch.renderFoodCard(ctx, b, mes, card, foodCards, "", lang)
}

func (fch *FoodCardHandler) renderFoodCard(
	ctx context.Context,
	b *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	card entities.FoodCard,
	foodCards []entities.FoodCard,
	info string,
	lang entities.LanguageCode,
) {
	msgText := fch.buildFoodText(lang, card, info)
	kb := fch.buildFoodKeyboard(mes.Message.Chat.ID, lang, card, foodCards)
	fch.sendFoodCardMessage(ctx, b, mes.Message.Chat.ID, card, msgText, kb)
	fch.deleteMessage(ctx, b, mes)
}

func (fch *FoodCardHandler) renderEmptyState(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	lang := fch.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
	m := utils.NewMessage(ctx, b, nil)
	kb := keyboards.GetEmptyFoodKeyboard(b, keyboards.EmptyFoodKeyboardConfig{
		AddFoodHandler: fch.addFoodButton,
		Lang:           lang,
	})
	m.Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyEmptyFoodCardText), true, kb)
	fch.deleteMessage(ctx, b, mes)
	m.SendSticker(mes.Message.Chat.ID)
}

func (fch *FoodCardHandler) buildFoodKeyboard(chatID int64, lang entities.LanguageCode, selectedCard entities.FoodCard, foodCards []entities.FoodCard) *inline.Keyboard {
	userID := entities.UserTelegramID(chatID)
	cartItems := fch.cartService.GetFoodCart(userID)

	buttons := make([]keyboards.FoodButton, 0, len(foodCards))
	for _, card := range foodCards {
		card := card
		buttons = append(buttons, keyboards.FoodButton{
			Text:    card.Name,
			OnClick: fch.GetOnSelectFood(card.Key),
		})
	}

	var showCart inline.OnSelect
	if fch.cartHandler != nil {
		showCart = fch.cartHandler.GetOnShow()
	}

	return keyboards.GetFoodKeyboard(fch.bot, keyboards.FoodKeyboardConfig{
		RandomHandler:        fch.GetOnSelectRandom(),
		AddToCart:            fch.addToCartHandler(selectedCard),
		ShowCartHandler:      showCart,
		MainMenuHandler:      fch.mainMenuFn,
		DeleteCurrentHandler: fch.deleteFoodHandler(selectedCard),
		CartButtonText:       fmt.Sprintf("%s (%d)", utils.Text(lang, utils.KeyBtnShowCart), len(cartItems)),
		AddFoodHandler:       fch.addFoodButton,
		EditCurrentHandler:   fch.getEditHandler(selectedCard),
		FoodButtons:          buttons,
	}, lang)
}

func (fch *FoodCardHandler) getEditHandler(card entities.FoodCard) inline.OnSelect {
	if fch.editFoodFunc == nil {
		return nil
	}
	return fch.editFoodFunc(card)
}

func (fch *FoodCardHandler) addToCartHandler(card entities.FoodCard) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer fch.recoverPanic()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		fch.cartService.AddFood(userID, card)
		foodCards := fch.foodCardService.GetAllFoodCards()
		lang := fch.langService.Get(userID)
		fch.renderFoodCard(ctx, b, mes, card, foodCards, utils.Text(lang, utils.KeyFoodAddedToCartText), lang)
	}
}

func (fch *FoodCardHandler) deleteFoodHandler(card entities.FoodCard) inline.OnSelect {
	if fch.editFoodFunc == nil {
		return nil
	}
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer fch.recoverPanic()
		if mes.Message == nil {
			return
		}
		if err := fch.foodCardService.DeleteFoodCard(card.Key); err != nil {
			lang := fch.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
			m := utils.NewMessage(ctx, b, nil)
			m.Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormCardNotFoundText), true, nil)
			return
		}
		lang := fch.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
		foodCards := fch.foodCardService.GetAllFoodCards()
		if len(foodCards) == 0 {
			fch.renderEmptyState(ctx, b, mes)
			return
		}
		nextCard := foodCards[rand.IntN(len(foodCards))]
		fch.renderFoodCard(ctx, b, mes, nextCard, foodCards, utils.Text(lang, utils.KeyFoodDeletedText), lang)
	}
}

func (fch *FoodCardHandler) findFoodCard(list []entities.FoodCard, key entities.FoodCardKey) (entities.FoodCard, bool) {
	for _, card := range list {
		if card.Key == key {
			return card, true
		}
	}
	return entities.FoodCard{}, false
}

func (fch *FoodCardHandler) buildFoodText(lang entities.LanguageCode, card entities.FoodCard, info string) string {
	nameText := utils.EscapeMarkdownV2(card.Name)
	descText := utils.EscapeMarkdownV2(card.Description)
	priceText := utils.EscapeMarkdownV2(utils.FormatUintSlice(card.Price))
	currencyText := utils.EscapeMarkdownV2(strings.Join(card.Currency, " "))
	timeText := utils.EscapeMarkdownV2(fmt.Sprintf("%d", card.TimeCooking))
	msgText := fmt.Sprintf(
		utils.Text(lang, utils.KeyFoodCardText),
		nameText,
		descText,
		priceText,
		currencyText,
		timeText,
	)
	if info != "" {
		msgText = fmt.Sprintf("%s\n\n_%s_", msgText, utils.EscapeMarkdownV2(info))
	}
	return msgText
}

func (fch *FoodCardHandler) sendFoodCardMessage(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	card entities.FoodCard,
	text string,
	kb *inline.Keyboard,
) {
	m := utils.NewMessage(ctx, b, nil)
	if card.PhotoFilePath != "" {
		photo, cleanup, err := fch.buildPhotoInput(card.PhotoFilePath)
		if err == nil {
			defer cleanup()
			if _, err := m.SendPhoto(chatID, photo, text, true, kb); err == nil {
				return
			}
		}
	}
	m.Send(chatID, text, true, kb)
}

func (fch *FoodCardHandler) buildPhotoInput(photoPath string) (models.InputFile, func(), error) {
	if photoPath == "" {
		return nil, func() {}, fmt.Errorf("empty photo path")
	}
	if info, err := os.Stat(photoPath); err == nil && !info.IsDir() {
		file, err := os.Open(photoPath)
		if err == nil {
			return &models.InputFileUpload{
				Filename: filepath.Base(photoPath),
				Data:     file,
			}, func() { _ = file.Close() }, nil
		}
	}
	return &models.InputFileString{Data: photoPath}, func() {}, nil
}

func (fch *FoodCardHandler) deleteMessage(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
	})
	if err != nil {
		log.Printf("delete message failed chat_id=%d message_id=%d err=%v", mes.Message.Chat.ID, mes.Message.ID, err)
	}
}

func (fch *FoodCardHandler) recoverPanic() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}
