package language

import "github.com/savabush/breakfastLoveBot/internal/entities"

type Repository interface {
	Get(userID entities.UserTelegramID) (entities.LanguageCode, bool)
	Set(userID entities.UserTelegramID, lang entities.LanguageCode) error
}
