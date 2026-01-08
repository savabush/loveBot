package middlewares

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/store"
)

func DeletePreviousBotMessage(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		if update.CallbackQuery != nil {
			next(ctx, b, update)
			return
		}
		if update.Message != nil {
			fmt.Println(update.Message.Text)
			messageStore := store.NewMessageStoreService(ctx, b)
			messageStore.Delete(entities.UserTelegramID(update.Message.Chat.ID))
		}
		next(ctx, b, update)
	}
}
