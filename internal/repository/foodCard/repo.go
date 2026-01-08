package foodCard

import "github.com/savabush/breakfastLoveBot/internal/entities"

type Repository interface {
	AddNewFood(foodCard entities.FoodCard)
	GetFoodCardByID(id entities.FoodCardKey) (entities.FoodCard, error)
	GetAllFoodCards() []entities.FoodCard
	UpdateFoodCard(foodCard entities.FoodCard) error
	DeleteFoodCard(id entities.FoodCardKey) error
}
