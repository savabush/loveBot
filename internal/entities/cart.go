package entities

type UserTelegramID int

type Cart struct {
	FoodItems map[UserTelegramID][]FoodCard
}
