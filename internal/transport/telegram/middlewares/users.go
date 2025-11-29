package middlewares

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/savabush/breakfastLoveBot/internal/config"
)

func CheckUserIDs(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		if update.Message != nil {
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			})
			if update.Message.Chat.ID != config.Cfg.UserID1 && update.Message.Chat.ID != config.Cfg.UserID2 {
				log.Println("who are u? ", update.Message.Chat.ID)
				return
			} else {
				next(ctx, b, update)
				return
			}
		}
		next(ctx, b, update)
	}
}
