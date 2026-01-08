package usecases

import (
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/sticker"
)

type StickerService struct {
	repo sticker.Repository
}

func NewStickerService(repo sticker.Repository) *StickerService {
	return &StickerService{repo: repo}
}

func (s *StickerService) AddSticker(sticker entities.Sticker) error {
	return s.repo.Add(sticker)
}

func (s *StickerService) GetNext() entities.Sticker {
	return s.repo.GetNext()
}

func (s *StickerService) HasStickers() bool {
	return s.repo.HasStickers()
}

func (s *StickerService) DeleteSticker(code string) error {
	return s.repo.Delete(code)
}
