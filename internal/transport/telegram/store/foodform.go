package store

import (
	"sync"

	"github.com/savabush/breakfastLoveBot/internal/entities"
)

type FoodFormStage int

const (
	FoodFormStageNone FoodFormStage = iota
	FoodFormStageName
	FoodFormStageDescription
	FoodFormStagePrice
	FoodFormStageCurrency
	FoodFormStageTime
	FoodFormStagePhoto
	FoodFormStageConfirm
)

type FoodFormMode int

const (
	FoodFormModeNone FoodFormMode = iota
	FoodFormModeAdd
	FoodFormModeEdit
)

type FoodFormSession struct {
	Mode        FoodFormMode
	Stage       FoodFormStage
	Draft       entities.FoodCard
	OriginalKey entities.FoodCardKey
}

type FoodFormStore struct {
	mu       sync.RWMutex
	sessions map[entities.UserTelegramID]*FoodFormSession
}

func NewFoodFormStore() *FoodFormStore {
	return &FoodFormStore{
		sessions: make(map[entities.UserTelegramID]*FoodFormSession),
	}
}

func (s *FoodFormStore) Save(userID entities.UserTelegramID, session *FoodFormSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[userID] = session
}

func (s *FoodFormStore) Get(userID entities.UserTelegramID) (*FoodFormSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[userID]
	return session, ok
}

func (s *FoodFormStore) Delete(userID entities.UserTelegramID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, userID)
}

func (s *FoodFormStore) Has(userID entities.UserTelegramID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.sessions[userID]
	return ok
}
