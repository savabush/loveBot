package usecases

import (
	"github.com/savabush/breakfastLoveBot/internal/entities"
	"github.com/savabush/breakfastLoveBot/internal/repository/cart"
)

type CartService struct {
	repo cart.Repository
}

func NewCartService(repo cart.Repository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddFood(userID entities.UserTelegramID, foodCard entities.FoodCard) {
	s.repo.AddFood(userID, foodCard)
}

func (s *CartService) GetFoodCart(userID entities.UserTelegramID) []entities.FoodCard {
	return s.repo.GetFoodCart(userID)
}

func (s *CartService) AcceptFoodCart(userID entities.UserTelegramID) {
	s.repo.AcceptFoodCart(userID)
}

func (s *CartService) CleanFoodCart(userID entities.UserTelegramID) {
	s.repo.CleanFoodCart(userID)
}
