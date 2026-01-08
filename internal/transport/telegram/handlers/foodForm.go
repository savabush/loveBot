package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/store"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type FoodFormHandler struct {
	bot             *bot.Bot
	foodCardService *usecases.FoodCardService
	store           *store.FoodFormStore
	stickerService  *usecases.StickerService
	langService     *usecases.LanguageService
	afterCancel     func(ctx context.Context, chatID int64)
	afterSave       func(ctx context.Context, chatID int64)
}

func NewFoodFormHandler(bot *bot.Bot, foodCardService *usecases.FoodCardService, stickerService *usecases.StickerService, langService *usecases.LanguageService) *FoodFormHandler {
	return &FoodFormHandler{
		bot:             bot,
		foodCardService: foodCardService,
		store:           store.NewFoodFormStore(),
		stickerService:  stickerService,
		langService:     langService,
	}
}

func (ffh *FoodFormHandler) SetAfterSave(fn func(ctx context.Context, chatID int64)) {
	ffh.afterSave = fn
}

func (ffh *FoodFormHandler) SetAfterCancel(fn func(ctx context.Context, chatID int64)) {
	ffh.afterCancel = fn
}

func (ffh *FoodFormHandler) Init() {
	ffh.bot.RegisterHandlerMatchFunc(ffh.matchActiveSession, ffh.handleActiveSession)
}

func (ffh *FoodFormHandler) matchActiveSession(update *models.Update) bool {
	if update.Message == nil {
		return false
	}
	userID := entities.UserTelegramID(update.Message.Chat.ID)
	session, ok := ffh.store.Get(userID)
	if !ok {
		return false
	}

	if len(update.Message.Photo) > 0 {
		return session.Stage == store.FoodFormStagePhoto
	}

	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return false
	}

	if strings.HasPrefix(text, "/") && !strings.EqualFold(text, "/cancel") {
		return false
	}

	return true
}

func (ffh *FoodFormHandler) handleActiveSession(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer ffh.recover()

	msg := update.Message
	userID := entities.UserTelegramID(msg.Chat.ID)
	session, ok := ffh.store.Get(userID)
	if !ok {
		return
	}

	text := strings.TrimSpace(msg.Text)
	if strings.EqualFold(text, "/cancel") || strings.EqualFold(text, "отмена") {
		ffh.cancelSession(ctx, userID, msg.Chat.ID)
		return
	}

	switch session.Stage {
	case store.FoodFormStageName:
		ffh.handleName(ctx, userID, msg.Chat.ID, session, text)
	case store.FoodFormStageDescription:
		ffh.handleDescription(ctx, userID, msg.Chat.ID, session, text)
	case store.FoodFormStagePrice:
		ffh.handlePrice(ctx, userID, msg.Chat.ID, session, text)
	case store.FoodFormStageCurrency:
		ffh.handleCurrency(ctx, userID, msg.Chat.ID, session, text)
	case store.FoodFormStageTime:
		ffh.handleTime(ctx, userID, msg.Chat.ID, session, text)
	case store.FoodFormStagePhoto:
		if len(msg.Photo) > 0 {
			ffh.handlePhotoFile(ctx, userID, msg.Chat.ID, session, msg.Photo)
		} else {
			ffh.handlePhoto(ctx, userID, msg.Chat.ID, session, text)
		}
	default:
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, msg.Chat.ID, utils.Text(lang, utils.KeyFoodFormUnknownInputText), true, nil)
	}
}

func (ffh *FoodFormHandler) GetOnAdd() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer ffh.recover()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		lang := ffh.langService.Get(userID)
		session := &store.FoodFormSession{
			Mode:  store.FoodFormModeAdd,
			Stage: store.FoodFormStageName,
		}
		ffh.store.Save(userID, session)
		ffh.sendMessage(ctx, mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormStartAddText), true, nil)
		ffh.promptNextStage(ctx, mes.Message.Chat.ID, session)
	}
}

func (ffh *FoodFormHandler) GetOnEdit(key entities.FoodCardKey) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer ffh.recover()
		card, err := ffh.foodCardService.GetFoodCardByID(key)
		if err != nil {
			lang := ffh.langService.Get(entities.UserTelegramID(mes.Message.Chat.ID))
			ffh.sendMessage(ctx, mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormCardNotFoundText), true, nil)
			return
		}
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		session := &store.FoodFormSession{
			Mode:        store.FoodFormModeEdit,
			Stage:       store.FoodFormStageConfirm,
			Draft:       card,
			OriginalKey: card.Key,
		}
		ffh.store.Save(userID, session)
		ffh.sendSummary(ctx, mes.Message.Chat.ID, session)
	}
}

func (ffh *FoodFormHandler) handleName(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	if value == "" {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormNamePrompt), true, nil)
		return
	}
	session.Draft.Name = value
	session.Stage = store.FoodFormStageDescription
	ffh.promptNextStage(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handleDescription(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	if value == "" {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormDescriptionPrompt), true, nil)
		return
	}
	session.Draft.Description = value
	session.Stage = store.FoodFormStagePrice
	ffh.promptNextStage(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handlePrice(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	priceItems, err := parseUintSlice(value)
	if err != nil || len(priceItems) == 0 {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormPricePrompt), true, nil)
		return
	}
	session.Draft.Price = priceItems
	session.Stage = store.FoodFormStageCurrency
	ffh.promptNextStage(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handleCurrency(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	if value == "" {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormCurrencyPrompt), true, nil)
		return
	}
	session.Draft.Currency = parseCurrencyList(value)
	session.Stage = store.FoodFormStageTime
	ffh.promptNextStage(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handleTime(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	minutes, err := strconv.Atoi(value)
	if err != nil || minutes < 0 {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormTimePrompt), true, nil)
		return
	}
	session.Draft.TimeCooking = uint(minutes)
	session.Stage = store.FoodFormStagePhoto
	ffh.promptNextStage(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handlePhoto(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, value string) {
	value = strings.TrimSpace(value)
	if value == "-" || value == "" {
		value = ""
	}
	session.Draft.PhotoFilePath = value
	session.Stage = store.FoodFormStageConfirm
	ffh.sendSummary(ctx, chatID, session)
}

func (ffh *FoodFormHandler) handlePhotoFile(ctx context.Context, userID entities.UserTelegramID, chatID int64, session *store.FoodFormSession, photos []models.PhotoSize) {
	if len(photos) == 0 {
		lang := ffh.langService.Get(userID)
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormPhotoPrompt), true, nil)
		return
	}
	best := photos[len(photos)-1]
	session.Draft.PhotoFilePath = best.FileID
	session.Stage = store.FoodFormStageConfirm
	ffh.sendSummary(ctx, chatID, session)
}

func (ffh *FoodFormHandler) promptNextStage(ctx context.Context, chatID int64, session *store.FoodFormSession) {
	lang := ffh.langService.Get(entities.UserTelegramID(chatID))
	switch session.Stage {
	case store.FoodFormStageName:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormNamePrompt), true, nil)
	case store.FoodFormStageDescription:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormDescriptionPrompt), true, nil)
	case store.FoodFormStagePrice:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormPricePrompt), true, nil)
	case store.FoodFormStageCurrency:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormCurrencyPrompt), true, nil)
	case store.FoodFormStageTime:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormTimePrompt), true, nil)
	case store.FoodFormStagePhoto:
		ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormPhotoPrompt), true, nil)
	case store.FoodFormStageConfirm:
		ffh.sendSummary(ctx, chatID, session)
	}
}

func (ffh *FoodFormHandler) sendSummary(ctx context.Context, chatID int64, session *store.FoodFormSession) {
	lang := ffh.langService.Get(entities.UserTelegramID(chatID))
	text := fmt.Sprintf(utils.Text(lang, utils.KeyFoodFormSummaryText),
		utils.EscapeMarkdownV2(session.Draft.Name),
		utils.EscapeMarkdownV2(session.Draft.Description),
		utils.EscapeMarkdownV2(utils.FormatUintSlice(session.Draft.Price)),
		utils.EscapeMarkdownV2(strings.Join(session.Draft.Currency, " ")),
		utils.EscapeMarkdownV2(fmt.Sprintf("%d", session.Draft.TimeCooking)),
		utils.EscapeMarkdownV2(displayPhoto(session.Draft.PhotoFilePath)),
	)

	var saveLabel string
	if session.Mode == store.FoodFormModeEdit {
		saveLabel = utils.Text(lang, utils.KeyBtnFoodFormUpdate)
	}

	kb := keyboards.GetFoodFormConfirmKeyboard(ffh.bot, keyboards.FoodFormKeyboardConfig{
		SaveHandler:         ffh.getSaveHandler(),
		CancelHandler:       ffh.getCancelHandler(),
		EditNameHandler:     ffh.getFieldEditHandler(store.FoodFormStageName),
		EditDescHandler:     ffh.getFieldEditHandler(store.FoodFormStageDescription),
		EditPriceHandler:    ffh.getFieldEditHandler(store.FoodFormStagePrice),
		EditCurrencyHandler: ffh.getFieldEditHandler(store.FoodFormStageCurrency),
		EditTimeHandler:     ffh.getFieldEditHandler(store.FoodFormStageTime),
		EditPhotoHandler:    ffh.getFieldEditHandler(store.FoodFormStagePhoto),
		SaveLabel:           saveLabel,
		Lang:                lang,
	})
	ffh.sendMessage(ctx, chatID, text, true, kb)
}

func (ffh *FoodFormHandler) getFieldEditHandler(target store.FoodFormStage) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer ffh.recover()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		session, ok := ffh.store.Get(userID)
		if !ok {
			return
		}
		session.Stage = target
		ffh.promptNextStage(ctx, mes.Message.Chat.ID, session)
	}
}

func (ffh *FoodFormHandler) getSaveHandler() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer ffh.recover()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		session, ok := ffh.store.Get(userID)
		if !ok {
			return
		}
		if session.Stage != store.FoodFormStageConfirm {
			return
		}

		if session.Mode == store.FoodFormModeAdd {
			if session.Draft.Key == "" {
				session.Draft.Key = entities.FoodCardKey(fmt.Sprintf("food-%d", rand.Uint64()))
			}
			ffh.foodCardService.AddNewFood(session.Draft)
			lang := ffh.langService.Get(userID)
			ffh.sendMessage(ctx, mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormSavedText), true, nil)
		} else if session.Mode == store.FoodFormModeEdit {
			session.Draft.Key = session.OriginalKey
			err := ffh.foodCardService.UpdateFoodCard(session.Draft)
			if err != nil {
				lang := ffh.langService.Get(userID)
				ffh.sendMessage(ctx, mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormCardNotFoundText), true, nil)
			} else {
				lang := ffh.langService.Get(userID)
				ffh.sendMessage(ctx, mes.Message.Chat.ID, utils.Text(lang, utils.KeyFoodFormUpdatedText), true, nil)
			}
		}
		ffh.store.Delete(userID)
		if ffh.afterSave != nil {
			ffh.afterSave(ctx, mes.Message.Chat.ID)
		}
	}
}

func (ffh *FoodFormHandler) getCancelHandler() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer ffh.recover()
		userID := entities.UserTelegramID(mes.Message.Chat.ID)
		ffh.cancelSession(ctx, userID, mes.Message.Chat.ID)
	}
}

func (ffh *FoodFormHandler) cancelSession(ctx context.Context, userID entities.UserTelegramID, chatID int64) {
	ffh.store.Delete(userID)
	lang := ffh.langService.Get(userID)
	ffh.sendMessage(ctx, chatID, utils.Text(lang, utils.KeyFoodFormCancelledText), true, nil)
	if ffh.afterCancel != nil {
		ffh.afterCancel(ctx, chatID)
	}
}

func (ffh *FoodFormHandler) sendMessage(ctx context.Context, chatID int64, text string, withDelete bool, kb *inline.Keyboard) {
	utils.NewMessage(ctx, ffh.bot, ffh.stickerService).Send(chatID, text, withDelete, kb)
}

func (ffh *FoodFormHandler) recover() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}

func parseUintSlice(raw string) ([]uint, error) {
	raw = strings.ReplaceAll(raw, ",", " ")
	parts := strings.Fields(raw)
	result := make([]uint, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if value < 0 {
			continue
		}
		result = append(result, uint(value))
	}
	return result, nil
}

func parseCurrencyList(raw string) entities.CurrencyList {
	raw = strings.ReplaceAll(raw, ",", " ")
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return entities.CurrencyList{}
	}
	return entities.CurrencyList(parts)
}

func displayPhoto(path string) string {
	if path == "" {
		return "—"
	}
	return "получено"
}
