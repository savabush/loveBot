package usecases

import (
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/language"
)

type LanguageService struct {
	repo language.Repository
}

func NewLanguageService(repo language.Repository) *LanguageService {
	return &LanguageService{repo: repo}
}

func (s *LanguageService) Get(userID entities.UserTelegramID) entities.LanguageCode {
	lang, ok := s.repo.Get(userID)
	if !ok || lang == "" {
		return entities.LanguageEN
	}
	return lang
}

func (s *LanguageService) Set(userID entities.UserTelegramID, lang entities.LanguageCode) {
	_ = s.repo.Set(userID, lang)
}

func (s *LanguageService) Toggle(userID entities.UserTelegramID) entities.LanguageCode {
	current := s.Get(userID)
	if current == entities.LanguageRU {
		s.Set(userID, entities.LanguageEN)
		return entities.LanguageEN
	}
	s.Set(userID, entities.LanguageRU)
	return entities.LanguageRU
}
