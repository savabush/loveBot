package store

import (
	"sync"

	"github.com/savabush/breakfastLoveBot/internal/entities"
)

type StickerFormSession struct{}

type StickerFormStore struct {
	mu       sync.RWMutex
	sessions map[entities.UserTelegramID]*StickerFormSession
}

func NewStickerFormStore() *StickerFormStore {
	return &StickerFormStore{
		sessions: make(map[entities.UserTelegramID]*StickerFormSession),
	}
}

func (s *StickerFormStore) Save(userID entities.UserTelegramID, session *StickerFormSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[userID] = session
}

func (s *StickerFormStore) Get(userID entities.UserTelegramID) (*StickerFormSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[userID]
	return session, ok
}

func (s *StickerFormStore) Delete(userID entities.UserTelegramID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, userID)
}

func (s *StickerFormStore) Has(userID entities.UserTelegramID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.sessions[userID]
	return ok
}
