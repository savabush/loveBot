package store

import (
	"context"
	"log"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/savabush/breakfastLoveBot/internal/entities"
)

var (
	MessageStore *MessageStoreService
	once         sync.Once
)

type MessageStoreService struct {
	b    *bot.Bot
	ctx  context.Context
	mu   *sync.RWMutex
	data map[entities.UserTelegramID][]*bot.DeleteMessageParams
}

func NewMessageStoreService(ctx context.Context, b *bot.Bot) *MessageStoreService {
	once.Do(func() {
		log.Println("create new message store")
		MessageStore = &MessageStoreService{
			b:    b,
			ctx:  ctx,
			mu:   &sync.RWMutex{},
			data: make(map[entities.UserTelegramID][]*bot.DeleteMessageParams),
		}
	})

	return MessageStore
}

func (s *MessageStoreService) Add(userID entities.UserTelegramID, params *bot.DeleteMessageParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[userID] = append(s.data[userID], params)
}

func (s *MessageStoreService) Delete(userID entities.UserTelegramID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, params := range s.data[userID] {
		s.b.DeleteMessage(s.ctx, params)
	}
	delete(s.data, userID)
}

func (s *MessageStoreService) Get(userID entities.UserTelegramID) []*bot.DeleteMessageParams {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[userID]
}
