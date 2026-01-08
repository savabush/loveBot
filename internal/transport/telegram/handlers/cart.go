package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/savabush/breakfastLoveBot/internal/config"
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/keyboards"
	"github.com/savabush/breakfastLoveBot/internal/transport/telegram/utils"
	"github.com/savabush/breakfastLoveBot/internal/usecases"
)

type CartHandler struct {
	bot          *bot.Bot
	cartService  *usecases.CartService
	backToFoodFn inline.OnSelect
	langService  *usecases.LanguageService
	mu           *sync.RWMutex
	pending      map[entities.UserTelegramID]cartOrder
	partial      map[entities.UserTelegramID]*partialSelection
	approvals    map[entities.UserTelegramID]approvalRequest
}

type cartOrder struct {
	From  entities.UserTelegramID
	Items []entities.FoodCard
}

type partialSelection struct {
	From     entities.UserTelegramID
	Items    []entities.FoodCard
	Selected map[int]bool
}

type approvalRequest struct {
	From  entities.UserTelegramID
	Items []entities.FoodCard
}

func NewCartHandler(bot *bot.Bot, cartService *usecases.CartService, langService *usecases.LanguageService) *CartHandler {
	return &CartHandler{
		bot:         bot,
		cartService: cartService,
		langService: langService,
		mu:          &sync.RWMutex{},
		pending:     make(map[entities.UserTelegramID]cartOrder),
		partial:     make(map[entities.UserTelegramID]*partialSelection),
		approvals:   make(map[entities.UserTelegramID]approvalRequest),
	}
}

func (ch *CartHandler) SetBackHandler(handler inline.OnSelect) {
	ch.backToFoodFn = handler
}

func (ch *CartHandler) GetOnShow() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.renderCart(ctx, bot, mes, "")
		})
	}
}

func (ch *CartHandler) GetOnClear() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			userID := entities.UserTelegramID(mes.Message.Chat.ID)
			ch.cartService.CleanFoodCart(userID)
			lang := ch.langService.Get(userID)
			ch.renderCart(ctx, bot, mes, utils.Text(lang, utils.KeyCartClearedText))
		})
	}
}

func (ch *CartHandler) GetOnAccept() inline.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			userID := entities.UserTelegramID(mes.Message.Chat.ID)
			items := ch.cartService.GetFoodCart(userID)
			ch.cartService.AcceptFoodCart(userID)
			ch.notifyOtherUser(ctx, bot, userID, items)
			lang := ch.langService.Get(userID)
			ch.renderCart(ctx, bot, mes, utils.Text(lang, utils.KeyCartAcceptedText))
		})
	}
}

func (ch *CartHandler) renderCart(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, info string) {
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	lang := ch.langService.Get(userID)
	items := ch.cartService.GetFoodCart(userID)

	if len(items) == 0 {
		text := utils.Text(lang, utils.KeyEmptyCartText)
		if info != "" {
			text = fmt.Sprintf("%s\n\n%s", info, text)
		}
		kb := keyboards.GetBackKeyboard(b, lang, ch.backToFoodFn)
		ch.respond(ctx, b, mes, text, kb)
		return
	}

	content := ch.formatCart(lang, items)
	if info != "" {
		content = fmt.Sprintf("%s\n\n_%s_", content, info)
	}
	kb := keyboards.GetCartKeyboard(b, lang, ch.backToFoodFn, ch.GetOnClear(), ch.GetOnAccept())
	ch.respond(ctx, b, mes, content, kb)
}

func (ch *CartHandler) formatCart(lang entities.LanguageCode, items []entities.FoodCard) string {
	var builder strings.Builder
	for i, card := range items {
		nameText := utils.EscapeMarkdownV2(card.Name)
		priceText := ch.formatCardPrice(card)
		details := fmt.Sprintf(utils.Text(lang, utils.KeyCartItemDetails), priceText, card.TimeCooking)
		detailsText := utils.EscapeMarkdownV2(details)
		fmt.Fprintf(&builder, "%d\\. %s\n%s\n", i+1, nameText, detailsText)
	}
	return fmt.Sprintf(utils.Text(lang, utils.KeyCartListText), builder.String())
}

func (ch *CartHandler) formatCardPrice(card entities.FoodCard) string {
	if len(card.Price) == 0 || len(card.Currency) == 0 {
		if len(card.Price) == 0 {
			return "-"
		}
	}
	parts := make([]string, 0, len(card.Price))
	for i, price := range card.Price {
		currency := ch.currencyForIndex(card, i)
		if currency == "" {
			parts = append(parts, fmt.Sprintf("%d", price))
			continue
		}
		parts = append(parts, fmt.Sprintf("%d %s", price, currency))
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}

func (ch *CartHandler) recover(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	fn()
}

func (ch *CartHandler) notifyOtherUser(ctx context.Context, b *bot.Bot, userID entities.UserTelegramID, items []entities.FoodCard) {
	if len(items) == 0 {
		return
	}
	otherID, ok := ch.getOtherUserID(userID)
	if !ok {
		return
	}
	lang := ch.langService.Get(otherID)
	foodList, totalTime, priceSummary := ch.formatCartSummary(items)
	text := fmt.Sprintf(utils.Text(lang, utils.KeyCartAcceptedForOtherText), foodList, totalTime, priceSummary)
	ch.savePendingOrder(otherID, cartOrder{
		From:  userID,
		Items: items,
	})
	kb := keyboards.GetOrderDecisionKeyboard(
		b,
		lang,
		ch.getOnOrderAccept(),
		ch.getOnOrderPartial(),
		ch.getOnOrderDecline(),
	)
	utils.NewMessage(ctx, b, nil).Send(int64(otherID), text, true, kb)
}

func (ch *CartHandler) getOtherUserID(userID entities.UserTelegramID) (entities.UserTelegramID, bool) {
	if int64(userID) == config.Cfg.UserID1 {
		return entities.UserTelegramID(config.Cfg.UserID2), true
	}
	if int64(userID) == config.Cfg.UserID2 {
		return entities.UserTelegramID(config.Cfg.UserID1), true
	}
	return 0, false
}

func (ch *CartHandler) formatCartSummary(items []entities.FoodCard) (string, uint, string) {
	var listBuilder strings.Builder
	totalTime := uint(0)
	maxTime := uint(0)
	currencyTotals := make(map[string]uint)

	for _, card := range items {
		name := utils.EscapeMarkdownV2(card.Name)
		listBuilder.WriteString("• ")
		listBuilder.WriteString(name)
		listBuilder.WriteString("\n")
		totalTime += card.TimeCooking
		if card.TimeCooking > maxTime {
			maxTime = card.TimeCooking
		}

		for i, price := range card.Price {
			currency := ch.currencyForIndex(card, i)
			if currency == "" {
				continue
			}
			currencyTotals[utils.EscapeMarkdownV2(currency)] += price
		}
	}

	keys := make([]string, 0, len(currencyTotals))
	for key := range currencyTotals {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var priceBuilder strings.Builder
	for i, key := range keys {
		if i > 0 {
			priceBuilder.WriteString("\n")
		}
		fmt.Fprintf(&priceBuilder, "• %d %s", currencyTotals[key], key)
	}

	if len(items) > 2 && maxTime > 0 {
		totalTime = uint(math.Ceil(float64(maxTime) * 1.5))
	}
	return strings.TrimRight(listBuilder.String(), "\n"), totalTime, priceBuilder.String()
}

func (ch *CartHandler) currencyForIndex(card entities.FoodCard, index int) string {
	if len(card.Currency) == 0 {
		return ""
	}
	if index >= 0 && index < len(card.Currency) {
		return card.Currency[index]
	}
	return card.Currency[len(card.Currency)-1]
}

func (ch *CartHandler) getOnOrderAccept() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handleOrderAccept(ctx, b, mes)
		})
	}
}

func (ch *CartHandler) getOnOrderPartial() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handleOrderPartial(ctx, b, mes)
		})
	}
}

func (ch *CartHandler) getOnOrderDecline() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handleOrderDecline(ctx, b, mes)
		})
	}
}

func (ch *CartHandler) handleOrderAccept(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	order, ok := ch.getPendingOrder(userID)
	if !ok {
		return
	}
	lang := ch.langService.Get(order.From)
	foodList, totalTime, priceSummary := ch.formatCartSummary(order.Items)
	text := fmt.Sprintf(utils.Text(lang, utils.KeyCartOrderAcceptedText), foodList, totalTime, priceSummary)
	utils.NewMessage(ctx, b, nil).Send(int64(order.From), text, true, nil)
	ch.clearPendingOrder(userID)
	ch.clearPartialSelection(userID)
	langCurrent := ch.langService.Get(userID)
	utils.NewMessage(ctx, b, nil).Send(mes.Message.Chat.ID, utils.Text(langCurrent, utils.KeyCartDecisionSentText), true, nil)
}

func (ch *CartHandler) handleOrderDecline(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	order, ok := ch.getPendingOrder(userID)
	if !ok {
		return
	}
	lang := ch.langService.Get(order.From)
	utils.NewMessage(ctx, b, nil).Send(int64(order.From), utils.Text(lang, utils.KeyCartOrderDeclinedText), true, nil)
	ch.clearPendingOrder(userID)
	ch.clearPartialSelection(userID)
	langCurrent := ch.langService.Get(userID)
	utils.NewMessage(ctx, b, nil).Send(mes.Message.Chat.ID, utils.Text(langCurrent, utils.KeyCartDecisionSentText), true, nil)
}

func (ch *CartHandler) handleOrderPartial(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	order, ok := ch.getPendingOrder(userID)
	if !ok {
		return
	}
	selection := &partialSelection{
		From:     order.From,
		Items:    order.Items,
		Selected: make(map[int]bool, len(order.Items)),
	}
	for i := range order.Items {
		selection.Selected[i] = true
	}
	ch.savePartialSelection(userID, selection)
	lang := ch.langService.Get(userID)
	text := ch.buildPartialSelectionText(lang, selection)
	kb := ch.buildPartialSelectionKeyboard(b, lang, selection)
	utils.NewMessage(ctx, b, nil).Edit(mes.Message.Chat.ID, mes.Message.ID, text, kb)
}

func (ch *CartHandler) handlePartialToggle(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, index int) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	selection, ok := ch.getPartialSelection(userID)
	if !ok || index < 0 || index >= len(selection.Items) {
		return
	}
	selection.Selected[index] = !selection.Selected[index]
	ch.savePartialSelection(userID, selection)
	lang := ch.langService.Get(userID)
	text := ch.buildPartialSelectionText(lang, selection)
	kb := ch.buildPartialSelectionKeyboard(b, lang, selection)
	utils.NewMessage(ctx, b, nil).Edit(mes.Message.Chat.ID, mes.Message.ID, text, kb)
}

func (ch *CartHandler) handlePartialSend(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	selection, ok := ch.getPartialSelection(userID)
	if !ok {
		return
	}
	selectedItems := make([]entities.FoodCard, 0, len(selection.Items))
	for i, item := range selection.Items {
		if selection.Selected[i] {
			selectedItems = append(selectedItems, item)
		}
	}
	if len(selectedItems) == 0 {
		lang := ch.langService.Get(userID)
		utils.NewMessage(ctx, b, nil).Edit(mes.Message.Chat.ID, mes.Message.ID, utils.Text(lang, utils.KeyCartPartialSendEmptyText), nil)
		return
	}
	lang := ch.langService.Get(selection.From)
	foodList, totalTime, priceSummary := ch.formatCartSummary(selectedItems)
	text := fmt.Sprintf(utils.Text(lang, utils.KeyCartPartialRequestText), foodList, totalTime, priceSummary)
	kb := keyboards.GetOrderApprovalKeyboard(b, lang, ch.getOnOrderApprove(), ch.getOnOrderCancel())
	ch.saveApproval(selection.From, approvalRequest{
		From:  userID,
		Items: selectedItems,
	})
	utils.NewMessage(ctx, b, nil).Send(int64(selection.From), text, true, kb)
	ch.clearPartialSelection(userID)
	langCurrent := ch.langService.Get(userID)
	utils.NewMessage(ctx, b, nil).Edit(mes.Message.Chat.ID, mes.Message.ID, utils.Text(langCurrent, utils.KeyCartPartialSentText), nil)
}

func (ch *CartHandler) handlePartialCancel(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	order, ok := ch.getPendingOrder(userID)
	if !ok {
		return
	}
	foodList, totalTime, priceSummary := ch.formatCartSummary(order.Items)
	lang := ch.langService.Get(userID)
	text := fmt.Sprintf(utils.Text(lang, utils.KeyCartAcceptedForOtherText), foodList, totalTime, priceSummary)
	kb := keyboards.GetOrderDecisionKeyboard(
		b,
		lang,
		ch.getOnOrderAccept(),
		ch.getOnOrderPartial(),
		ch.getOnOrderDecline(),
	)
	utils.NewMessage(ctx, b, nil).Edit(mes.Message.Chat.ID, mes.Message.ID, text, kb)
	ch.clearPartialSelection(userID)
}

func (ch *CartHandler) getOnOrderApprove() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handleOrderApprove(ctx, b, mes)
		})
	}
}

func (ch *CartHandler) getOnOrderCancel() inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handleOrderCancel(ctx, b, mes)
		})
	}
}

func (ch *CartHandler) handleOrderApprove(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	approval, ok := ch.getApproval(userID)
	if !ok {
		return
	}
	lang := ch.langService.Get(approval.From)
	foodList, totalTime, priceSummary := ch.formatCartSummary(approval.Items)
	text := fmt.Sprintf(utils.Text(lang, utils.KeyCartPartialApprovedText), foodList, totalTime, priceSummary)
	utils.NewMessage(ctx, b, nil).Send(int64(approval.From), text, true, nil)
	ch.clearPendingOrder(approval.From)
	ch.clearApproval(userID)
	langCurrent := ch.langService.Get(userID)
	utils.NewMessage(ctx, b, nil).Send(mes.Message.Chat.ID, utils.Text(langCurrent, utils.KeyCartDecisionSentText), true, nil)
}

func (ch *CartHandler) handleOrderCancel(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	userID := entities.UserTelegramID(mes.Message.Chat.ID)
	approval, ok := ch.getApproval(userID)
	if !ok {
		return
	}
	lang := ch.langService.Get(approval.From)
	utils.NewMessage(ctx, b, nil).Send(int64(approval.From), utils.Text(lang, utils.KeyCartPartialDeclinedText), true, nil)
	ch.clearPendingOrder(approval.From)
	ch.clearApproval(userID)
	langCurrent := ch.langService.Get(userID)
	utils.NewMessage(ctx, b, nil).Send(mes.Message.Chat.ID, utils.Text(langCurrent, utils.KeyCartDecisionSentText), true, nil)
}

func (ch *CartHandler) buildPartialSelectionText(lang entities.LanguageCode, selection *partialSelection) string {
	var builder strings.Builder
	builder.WriteString(utils.Text(lang, utils.KeyCartPartialSelectionText))
	builder.WriteString("\n")
	for i, item := range selection.Items {
		marker := "▫️"
		if selection.Selected[i] {
			marker = "✅"
		}
		fmt.Fprintf(&builder, "%s %d\\. %s\n", marker, i+1, utils.EscapeMarkdownV2(item.Name))
	}
	return strings.TrimRight(builder.String(), "\n")
}

func (ch *CartHandler) buildPartialSelectionKeyboard(b *bot.Bot, lang entities.LanguageCode, selection *partialSelection) *inline.Keyboard {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	for i, item := range selection.Items {
		index := i
		marker := "▫️"
		if selection.Selected[i] {
			marker = "✅"
		}
		label := fmt.Sprintf("%s %d. %s", marker, i+1, item.Name)
		kb.Row().Button(label, []byte(fmt.Sprintf("partial-toggle-%d", i)), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			ch.recover(func() {
				ch.handlePartialToggle(ctx, b, mes, index)
			})
		})
	}
	kb.Row().Button(utils.Text(lang, utils.KeyBtnPartialSend), []byte("partial-send"), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handlePartialSend(ctx, b, mes)
		})
	})
	kb.Row().Button(utils.Text(lang, utils.KeyBtnPartialCancel), []byte("partial-cancel"), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		ch.recover(func() {
			ch.handlePartialCancel(ctx, b, mes)
		})
	})
	return kb
}

func (ch *CartHandler) savePendingOrder(userID entities.UserTelegramID, order cartOrder) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.pending[userID] = order
}

func (ch *CartHandler) getPendingOrder(userID entities.UserTelegramID) (cartOrder, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	order, ok := ch.pending[userID]
	return order, ok
}

func (ch *CartHandler) clearPendingOrder(userID entities.UserTelegramID) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	delete(ch.pending, userID)
}

func (ch *CartHandler) savePartialSelection(userID entities.UserTelegramID, selection *partialSelection) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.partial[userID] = selection
}

func (ch *CartHandler) getPartialSelection(userID entities.UserTelegramID) (*partialSelection, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	selection, ok := ch.partial[userID]
	return selection, ok
}

func (ch *CartHandler) clearPartialSelection(userID entities.UserTelegramID) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	delete(ch.partial, userID)
}

func (ch *CartHandler) saveApproval(userID entities.UserTelegramID, approval approvalRequest) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.approvals[userID] = approval
}

func (ch *CartHandler) getApproval(userID entities.UserTelegramID) (approvalRequest, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	approval, ok := ch.approvals[userID]
	return approval, ok
}

func (ch *CartHandler) clearApproval(userID entities.UserTelegramID) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	delete(ch.approvals, userID)
}

func (ch *CartHandler) respond(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, text string, kb *inline.Keyboard) {
	m := utils.NewMessage(ctx, b, nil)
	if mes.Message == nil {
		return
	}
	if mes.Message.Text == "" {
		m.Send(mes.Message.Chat.ID, text, true, kb)
		ch.deleteMessage(ctx, b, mes)
		return
	}
	m.Edit(mes.Message.Chat.ID, mes.Message.ID, text, kb)
}

func (ch *CartHandler) deleteMessage(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	if mes.Message == nil {
		return
	}
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
	})
	if err != nil {
		log.Printf("delete message failed chat_id=%d message_id=%d err=%v", mes.Message.Chat.ID, mes.Message.ID, err)
	}
}
