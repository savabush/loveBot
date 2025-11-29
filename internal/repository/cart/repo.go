package cart

import "github.com/savabush/breakfastLoveBot/internal/entities"

type Repository interface {
	AddFood(userID entities.UserTelegramID, foodCard entities.FoodCard)
	AcceptFoodCart(userID entities.UserTelegramID)
	GetFoodCart(userID entities.UserTelegramID) []entities.FoodCard
	CleanFoodCart(userID entities.UserTelegramID)
}
