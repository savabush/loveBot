package utils

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/sticker"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/store"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type Message struct {
	b   *bot.Bot
	ctx context.Context

	messageStore   *store.MessageStoreService
	stickerService *usecases.StickerService
}

func NewMessage(ctx context.Context, b *bot.Bot) *Message {
	stickerRepo := sticker.NewMemoryRepository()
	stickerService := usecases.NewStickerService(
		stickerRepo,
	)
	return &Message{
		b:              b,
		ctx:            ctx,
		messageStore:   store.NewMessageStoreService(ctx, b),
		stickerService: stickerService,
	}
}

func (m *Message) Send(chatID int64, text string, withDelete bool, kb *inline.Keyboard) {
	params := &bot.SendMessageParams{
		Text:      text,
		ChatID:    chatID,
		ParseMode: models.ParseModeMarkdown,
	}
	if kb != nil {
		params.ReplyMarkup = kb
	}
	msg, err := m.b.SendMessage(
		m.ctx,
		params,
	)
	if err != nil {
		log.Println("err to send message: ", err)
	}
	if withDelete {
		m.messageStore.Add(entities.UserTelegramID(msg.Chat.ID), &bot.DeleteMessageParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
		})
	}

}

func (m *Message) Edit(chatID int64, messageID int, text string, kb *inline.Keyboard) {
	params := &bot.EditMessageTextParams{
		Text:      text,
		ChatID:    chatID,
		ParseMode: models.ParseModeMarkdown,
	}
	if kb != nil {
		params.ReplyMarkup = kb
	}
	_, err := m.b.EditMessageText(
		m.ctx,
		params,
	)
	if err != nil {
		log.Println("err to send message: ", err)
	}
}

func (m *Message) SendSticker(chatID int64) {
	if !m.stickerService.HasStickers() {
		return
	}
	msg, err := m.b.SendSticker(m.ctx, &bot.SendStickerParams{
		ChatID: chatID,
		Sticker: &models.InputFileString{
			Data: m.stickerService.GetNext().Code,
		},
	})
	if err != nil {
		log.Println(err)
	}
	m.messageStore.Add(entities.UserTelegramID(msg.Chat.ID), &bot.DeleteMessageParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
	})
}

func (m *Message) HasStickers() bool {
	return m.stickerService.HasStickers()
}
