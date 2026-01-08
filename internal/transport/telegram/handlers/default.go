package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

func NewDefaultHandler(langService *usecases.LanguageService) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}
		userID := entities.UserTelegramID(update.Message.Chat.ID)
		lang := langService.Get(userID)
		m := utils.NewMessage(ctx, b, nil)
		if m.HasStickers() {
			m.Send(update.Message.Chat.ID, utils.Text(lang, utils.KeyDefaultMessageTextWithStickers), true, nil)
		} else {
			m.Send(update.Message.Chat.ID, utils.Text(lang, utils.KeyDefaultMessageText), true, nil)
		}
		m.SendSticker(update.Message.Chat.ID)
	}
}
