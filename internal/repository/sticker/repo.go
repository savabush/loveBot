package sticker

import "github.com/savabush/breakfastLoveBot/internal/entities"

type Repository interface {
	Add(sticker entities.Sticker) error
	GetNext() entities.Sticker
	HasStickers() bool
}
