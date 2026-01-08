package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type StartHandler struct {
	bot         *bot.Bot
	keyboardCfg keyboards.StartKeyboardConfig
	langService *usecases.LanguageService
}

func NewStartHandler(bot *bot.Bot, keyboardCfg keyboards.StartKeyboardConfig, langService *usecases.LanguageService) *StartHandler {
	return &StartHandler{
		bot:         bot,
		keyboardCfg: keyboardCfg,
		langService: langService,
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
	lang := sh.langService.Get(entities.UserTelegramID(msg.Chat.ID))
	name := strings.TrimSpace(strings.TrimSpace(msg.Chat.FirstName) + " " + strings.TrimSpace(msg.Chat.LastName))
	now := time.Unix(int64(msg.Date), 0).Hour()
	stageOfDay := stageOfDayLabel(lang, now)

	msgText := fmt.Sprintf(utils.Text(lang, utils.KeyStartMessageText), stageOfDay, utils.EscapeMarkdownV2(name))

	kb := keyboards.GetStartKeyboard(b, sh.keyboardCfg, lang)
	m := utils.NewMessage(ctx, b, nil)
	m.Send(msg.Chat.ID, msgText, true, kb)
	m.SendSticker(msg.Chat.ID)
}

func (sh *StartHandler) GetOnSelect() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()

		if mes.Message == nil {
			return
		}
		msg := mes.Message
		lang := sh.langService.Get(entities.UserTelegramID(msg.Chat.ID))
		name := strings.TrimSpace(strings.TrimSpace(msg.Chat.FirstName) + " " + strings.TrimSpace(msg.Chat.LastName))
		now := time.Now().Hour()
		stageOfDay := stageOfDayLabel(lang, now)
		msgText := fmt.Sprintf(utils.Text(lang, utils.KeyStartMessageText), stageOfDay, utils.EscapeMarkdownV2(name))
		kb := keyboards.GetStartKeyboard(b, sh.keyboardCfg, lang)
		m := utils.NewMessage(ctx, b, nil)
		m.Send(msg.Chat.ID, msgText, true, kb)
		m.SendSticker(msg.Chat.ID)
	}
}

func (sh *StartHandler) GetOnSwitchLanguage() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		if mes.Message == nil {
			return
		}
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		lang := sh.langService.Toggle(userID)
		name := strings.TrimSpace(strings.TrimSpace(mes.Message.Chat.FirstName) + " " + strings.TrimSpace(mes.Message.Chat.LastName))
		now := time.Now().Hour()
		stageOfDay := stageOfDayLabel(lang, now)
		msgText := fmt.Sprintf(utils.Text(lang, utils.KeyStartMessageText), stageOfDay, utils.EscapeMarkdownV2(name))
		kb := keyboards.GetStartKeyboard(b, sh.keyboardCfg, lang)
		utils.NewMessage(ctx, b, nil).Send(mes.Message.Chat.ID, msgText, true, kb)
	}
}

func stageOfDayLabel(lang entities.LanguageCode, hour int) string {
	if lang == entities.LanguageRU {
		switch {
		case hour >= 6 && hour < 12:
			return "Доброе утро"
		case hour >= 12 && hour < 18:
			return "Добрый день"
		case hour >= 18 && hour < 22:
			return "Добрый вечер"
		default:
			return "Доброй ночи"
		}
	}
	switch {
	case hour >= 6 && hour < 12:
		return "Good morning"
	case hour >= 12 && hour < 18:
		return "Good afternoon"
	case hour >= 18 && hour < 22:
		return "Good evening"
	default:
		return "Good night"
	}
}

func (sh *StartHandler) SetLanguageHandler(handler inline.OnSelect) {
	sh.keyboardCfg.LanguageHandler = handler
}

func (sh *StartHandler) Init() {
	sh.bot.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, sh.startHandlerAny)
}
