package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type FoodCardHandler struct {
	bot *bot.Bot

	foodCardService *usecases.FoodCardService
}

func NewFoodCardHandler(bot *bot.Bot, foodCardService *usecases.FoodCardService) *FoodCardHandler {
	return &FoodCardHandler{
		bot:             bot,
		foodCardService: foodCardService,
	}
}

func (sh *FoodCardHandler) foodCardHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	foodCards := sh.foodCardService.GetAllFoodCards()
	m := utils.NewMessage(ctx, b)
	if len(foodCards) < 1 {
		m.Edit(mes.Message.Chat.ID, mes.Message.ID, utils.EmptyFoodCardText, keyboards.GetStartKeyboard(b, sh.GetOnSelect()))
		m.SendSticker(mes.Message.Chat.ID)
		return
	}
	randomFoodCard := foodCards[rand.IntN(len(foodCards)-1)]
	msgText := fmt.Sprintf(
		utils.FoodCardText,
		randomFoodCard.Name,
		randomFoodCard.Description,
		randomFoodCard.Price,
		randomFoodCard.Currency,
		randomFoodCard.TimeCooking,
	)

	kb := keyboards.GetFoodKeyboard(b, 0)
	m.Edit(mes.Message.Chat.ID, mes.Message.ID, msgText, kb)
}

func (fch *FoodCardHandler) GetOnSelect() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		fch.foodCardHandler(ctx, bot, mes, data)
	}
}
