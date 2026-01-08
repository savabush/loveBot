package handlers

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/store"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type StickerHandler struct {
	bot            *bot.Bot
	stickerService *usecases.StickerService
	langService    *usecases.LanguageService
	store          *store.StickerFormStore
	mu             *sync.RWMutex
	lastShown      map[entities.UserTelegramID]string
	mainMenuFn     inline.OnSelect
}

func NewStickerHandler(bot *bot.Bot, stickerService *usecases.StickerService, langService *usecases.LanguageService) *StickerHandler {
	return &StickerHandler{
		bot:            bot,
		stickerService: stickerService,
		langService:    langService,
		store:          store.NewStickerFormStore(),
		mu:             &sync.RWMutex{},
		lastShown:      make(map[entities.UserTelegramID]string),
	}
}

func (sh *StickerHandler) Init() {
	sh.bot.RegisterHandlerMatchFunc(sh.matchActiveSticker, sh.handleSticker)
}

func (sh *StickerHandler) SetMainMenuHandler(handler inline.OnSelect) {
	sh.mainMenuFn = handler
}

func (sh *StickerHandler) GetOnSelect() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		sh.showSticker(ctx, b, mes)
	}
}

func (sh *StickerHandler) GetOnAdd() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		lang := sh.langService.Get(userID)
		sh.store.Save(userID, &store.StickerFormSession{})
		utils.NewMessage(ctx, b, sh.stickerService).
			Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyStickerFormPrompt), true, nil)
	}
}

func (sh *StickerHandler) matchActiveSticker(update *models.Update) bool {
	if update.Message == nil || update.Message.Sticker == nil {
		return false
	}
	userID := entities.UserTelegramID(update.Message.Chat.ID)
	return sh.store.Has(userID)
}

func (sh *StickerHandler) handleSticker(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	userID := entities.UserTelegramID(update.Message.Chat.ID)
	session, ok := sh.store.Get(userID)
	if !ok || session == nil {
		return
	}
	fileID := strings.TrimSpace(update.Message.Sticker.FileID)
	if fileID == "" {
		return
	}
	_ = sh.stickerService.AddSticker(entities.Sticker{Code: fileID})
	sh.store.Delete(userID)
	lang := sh.langService.Get(userID)
	utils.NewMessage(ctx, b, sh.stickerService).
		Send(update.Message.Chat.ID, utils.Text(lang, utils.KeyStickerFormSavedText), true, nil)
	sh.showSticker(ctx, b, models.MaybeInaccessibleMessage{Message: update.Message})
}

func (sh *StickerHandler) GetOnDelete() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		if mes.Message == nil {
			return
		}
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		lang := sh.langService.Get(userID)
		sh.mu.RLock()
		code := strings.TrimSpace(sh.lastShown[userID])
		sh.mu.RUnlock()
		if code == "" {
			utils.NewMessage(ctx, b, sh.stickerService).
				Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyStickerDeleteEmpty), true, nil)
			return
		}
		if err := sh.stickerService.DeleteSticker(code); err != nil {
			utils.NewMessage(ctx, b, sh.stickerService).
				Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyStickerDeleteEmpty), true, nil)
			return
		}
		sh.mu.Lock()
		sh.lastShown[userID] = ""
		sh.mu.Unlock()
		utils.NewMessage(ctx, b, sh.stickerService).
			Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyStickerDeletedText), true, nil)
		sh.showSticker(ctx, b, mes)
	}
}

func (sh *StickerHandler) showSticker(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	m := utils.NewMessage(ctx, b, sh.stickerService)
	lang := sh.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
	kb := keyboards.GetStickerKeyboard(b, keyboards.StickerKeyboardConfig{
		NextHandler:     sh.GetOnSelect(),
		AddHandler:      sh.GetOnAdd(),
		DeleteHandler:   sh.GetOnDelete(),
		MainMenuHandler: sh.mainMenuFn,
		Lang:            lang,
	})
	if !m.HasStickers() {
		m.Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyEmptyStickersText), true, kb)
		return
	}
	sticker := sh.stickerService.GetNext()
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	sh.mu.Lock()
	sh.lastShown[userID] = sticker.Code
	sh.mu.Unlock()
	m.Send(mes.Message.Chat.ID, utils.Text(lang, utils.KeyStickerCatchText), true, kb)
	_ = m.SendStickerWithCode(mes.Message.Chat.ID, sticker.Code)
}
