package utils

import (
	"context"
	"log"
	"strings"

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

func NewMessage(ctx context.Context, b *bot.Bot, stickerService *usecases.StickerService) *Message {
	if stickerService == nil {
		stickerRepo := sticker.NewMemoryRepository()
		stickerService = usecases.NewStickerService(
			stickerRepo,
		)
	}
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

func (m *Message) SendPhoto(chatID int64, photo models.InputFile, caption string, withDelete bool, kb *inline.Keyboard) (*models.Message, error) {
	params := &bot.SendPhotoParams{
		ChatID:    chatID,
		Photo:     photo,
		Caption:   caption,
		ParseMode: models.ParseModeMarkdown,
	}
	if kb != nil {
		params.ReplyMarkup = kb
	}
	msg, err := m.b.SendPhoto(m.ctx, params)
	if err != nil {
		log.Println("err to send photo: ", err)
		return nil, err
	}
	if withDelete {
		m.messageStore.Add(entities.UserTelegramID(msg.Chat.ID), &bot.DeleteMessageParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
		})
	}
	return msg, nil
}

func (m *Message) Edit(chatID int64, messageID int, text string, kb *inline.Keyboard) {
	params := &bot.EditMessageTextParams{
		Text:      text,
		ChatID:    chatID,
		MessageID: messageID,
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
		log.Printf("edit message failed chat_id=%d message_id=%d err=%v", chatID, messageID, err)
		if strings.Contains(err.Error(), "message to edit not found") {
			log.Printf("fallback to send message chat_id=%d message_id=%d", chatID, messageID)
			m.Send(chatID, text, true, kb)
		}
	}
}

func (m *Message) SendSticker(chatID int64) {
	if !m.stickerService.HasStickers() {
		return
	}
	_ = m.SendStickerWithCode(chatID, m.stickerService.GetNext().Code)
}

func (m *Message) SendStickerWithCode(chatID int64, code string) error {
	if strings.TrimSpace(code) == "" {
		return nil
	}
	msg, err := m.b.SendSticker(m.ctx, &bot.SendStickerParams{
		ChatID: chatID,
		Sticker: &models.InputFileString{
			Data: code,
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	m.messageStore.Add(entities.UserTelegramID(msg.Chat.ID), &bot.DeleteMessageParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
	})
	return nil
}

func (m *Message) HasStickers() bool {
	return m.stickerService.HasStickers()
}
