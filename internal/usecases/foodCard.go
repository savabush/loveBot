package usecases

import (
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/foodCard"
)

type FoodCardService struct {
	repo foodCard.Repository
}

func NewFoodCardService(repo foodCard.Repository) *FoodCardService {
	return &FoodCardService{repo: repo}
}

func (s *FoodCardService) AddNewFood(foodCard entities.FoodCard) {
	s.repo.AddNewFood(foodCard)
}

func (s *FoodCardService) GetFoodCardByID(id entities.FoodCardKey) (entities.FoodCard, error) {
	return s.repo.GetFoodCardByID(id)
}

func (s *FoodCardService) GetAllFoodCards() []entities.FoodCard {
	return s.repo.GetAllFoodCards()
}

func (s *FoodCardService) UpdateFoodCard(foodCard entities.FoodCard) error {
	return s.repo.UpdateFoodCard(foodCard)
}

func (s *FoodCardService) DeleteFoodCard(id entities.FoodCardKey) error {
	return s.repo.DeleteFoodCard(id)
}
