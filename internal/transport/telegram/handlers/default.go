package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	m := utils.NewMessage(ctx, b)
	if m.HasStickers() {
		m.Send(update.Message.Chat.ID, utils.DefaultMessageTextWithStickers, true, nil)
	} else {
		m.Send(update.Message.Chat.ID, utils.DefaultMessageText, true, nil)
	}
	m.SendSticker(update.Message.Chat.ID)
}
