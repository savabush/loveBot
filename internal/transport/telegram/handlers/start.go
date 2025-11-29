package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/savabush/breakfastLoveBot/internal/repository/foodCard"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type StartHandler struct {
	bot *bot.Bot
}

func NewStartHandler(bot *bot.Bot) *StartHandler {
	return &StartHandler{
		bot: bot,
	}
}

func (sh *StartHandler) startHandlerAny(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	var msg *models.Message
	msg = update.Message
	name := msg.Chat.FirstName + " " + msg.Chat.LastName
	now := time.Unix(int64(msg.Date), 0).Hour()
	var stageOfDay string
	switch {
	case now >= 6 && now < 12:
		stageOfDay = "Доброе утро"
	case now >= 12 && now < 18:
		stageOfDay = "Добрый день"
	case now >= 18 && now < 22:
		stageOfDay = "Добрый вечер"
	case (now >= 22 && now <= 24) || (now >= 0 && now < 6):
		stageOfDay = "Доброй ночи"
	}

	msgText := fmt.Sprintf(utils.StartMessageText, stageOfDay, name)

	foodCardRepo := foodCard.NewMemoryRepository()
	foodCardService := usecases.NewFoodCardService(foodCardRepo)
	foodCardHandler := NewFoodCardHandler(b, foodCardService)
	kb := keyboards.GetStartKeyboard(b, foodCardHandler.GetOnSelect())
	m := utils.NewMessage(ctx, b)
	m.Send(msg.Chat.ID, msgText, true, kb)
	m.SendSticker(msg.Chat.ID)
}

func (sh *StartHandler) Init() {
	sh.bot.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, sh.startHandlerAny)
}
